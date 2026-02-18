package handlers

import (
	"errors"
	"net/http"
	"time"

	"rest_api/internal/repository"
)

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	// 1) Достаём refresh token из cookie
	c, err := r.Cookie("refresh_token")
	if err != nil || c.Value == "" {
		respondError(w, http.StatusUnauthorized, "Необходим refresh token")
		return
	}

	refreshToken := c.Value

	// 2) Проверяем refresh JWT (подпись, exp, typ=refresh) и достаём userID
	userID, err := h.jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		clearRefreshCookie(w)
		respondError(w, http.StatusUnauthorized, "Неверный refresh token")
		return
	}

	// 3) Хешируем токен и ищем сессию в БД
	tokenHash := hashToken(refreshToken)

	session, err := h.refreshStore.GetByHash(r.Context(), tokenHash)
	if err != nil {
		clearRefreshCookie(w)
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusUnauthorized, "refresh session not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to load refresh session")
		return
	}

	// 4) Дополнительные проверки сессии
	// (JWT уже проверил exp токена, но мы проверяем expires_at на стороне сервера тоже)
	if session.UserID != userID {
		clearRefreshCookie(w)
		respondError(w, http.StatusUnauthorized, "refresh session mismatch")
		return
	}

	if session.RevokedAt != nil {
		clearRefreshCookie(w)
		respondError(w, http.StatusUnauthorized, "refresh session revoked")
		return
	}

	if time.Now().After(session.ExpiresAt) {
		clearRefreshCookie(w)
		respondError(w, http.StatusUnauthorized, "refresh session expired")
		return
	}

	// 5) ROTATION: сначала revoke старую сессию
	if err := h.refreshStore.RevokeByID(r.Context(), session.ID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to revoke refresh session")
		return
	}

	// 6) Генерим новый refresh token и сохраняем новую сессию
	newRefreshToken, err := h.jwt.GenerateRefreshToken(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}

	refreshTTL := h.jwt.RefreshTTL()
	newRefreshExpiresAt := time.Now().Add(refreshTTL)
	newHash := hashToken(newRefreshToken)

	if _, err := h.refreshStore.Create(r.Context(), userID, newHash, newRefreshExpiresAt); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to persist refresh session")
		return
	}

	// 7) Обновляем cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/auth/refresh",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,

		MaxAge:  int(refreshTTL.Seconds()),
		Expires: newRefreshExpiresAt,
	})

	// 8) Выдаём новый access token
	accessToken, err := h.jwt.GenerateAccessToken(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate access token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"access_token": accessToken,
	})
}

// маленький helper — чтобы не дублировать в нескольких местах
func clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth/refresh",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
