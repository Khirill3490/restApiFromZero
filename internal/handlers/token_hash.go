package handlers

import (
	"crypto/sha256"
	"encoding/hex"
)

// hashToken возвращает SHA-256 hash токена в hex-виде.
// Мы храним в БД именно hash, а не сырой refresh JWT.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
