package handlers

import (
	"net/http"
)

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("refresh_token")
	if err != nil {
		respondError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	userID, err := h.jwt.ParseRefreshToken(c.Value)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	accessToken, err := h.jwt.GenerateAccessToken(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate access token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"access_token": accessToken,
	})
}
