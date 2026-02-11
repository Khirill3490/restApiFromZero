package repository

import (
	"rest_api/internal/models"
	"time"
)


type RefreshSessionRepository interface {
	Create(userID int64, tokenHash string, expiresAt time.Time) (models.RefreshSession, error)
	GetByHash(tokenHash string) (models.RefreshSession, error)
	RevokeByHash(tokenHash string) error
}
