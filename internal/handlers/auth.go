package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

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
