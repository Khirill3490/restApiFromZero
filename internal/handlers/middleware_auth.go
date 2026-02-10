package handlers

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const ctxUserIDKey ctxKey = "userID"

func (h *Handler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			respondError(w, http.StatusUnauthorized, "Заголовок Authorization не найден")
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			respondError(w, http.StatusUnauthorized, "Неверный формат заголовка Authorization")
			return
		}

		userID, err := h.jwt.ParseAccessToken(parts[1])
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Неверный токен")
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	v := ctx.Value(ctxUserIDKey)
	id, ok := v.(int64)
	return id, ok
}
