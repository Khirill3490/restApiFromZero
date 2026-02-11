package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"rest_api/internal/models"
	"rest_api/internal/repository"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterInput

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}

	input.Email = strings.TrimSpace(input.Email)
	if input.Email == "" || input.Password == "" {
		respondError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	if len(input.Password) < 6 {
		respondError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user, err := h.userStore.Create(
		r.Context(),
		input.Email,
		string(passwordHash),
	)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			respondError(w, http.StatusConflict, "email already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.LoginInput

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "нкорректный json")
		return
	}

	input.Email = strings.TrimSpace(input.Email)
	if input.Email == "" || input.Password == "" {
		respondError(w, http.StatusBadRequest, "email и password обязательны")
		return
	}

	user, err := h.userStore.GetByEmail(r.Context(), input.Email)
	if err != nil {
		// Важно: чтобы не “палить”, существует ли email — возвращаем 401 одинаково.
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusUnauthorized, "неверные credentials	")
			return
		}
		respondError(w, http.StatusInternalServerError, "ошибка при получении пользователя")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		respondError(w, http.StatusUnauthorized, "неверные credentials")
		return
	}

	refreshToken, err := h.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "ошибка при генерации refresh токена")
		return
	}

	refreshTTL := h.jwt.RefreshTTL()

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/auth/refresh",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,

		MaxAge:  int(refreshTTL.Seconds()),
		Expires: time.Now().Add(refreshTTL),
	})

	accessToken, err := h.jwt.GenerateAccessToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "ошибка при генерации токена")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"access_token": accessToken,
		"user":         user,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth/refresh",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // удалить cookie
	})

	w.WriteHeader(http.StatusNoContent)
}
