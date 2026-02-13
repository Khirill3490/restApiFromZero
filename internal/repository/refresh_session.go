package repository

import (
	"context"
	"rest_api/internal/models"
	"time"
)

type RefreshSessionRepository interface {
	Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) (*models.RefreshSession, error)
	GetByHash(ctx context.Context, tokenHash string) (*models.RefreshSession, error)
	RevokeByID(ctx context.Context, id int64) error
}
