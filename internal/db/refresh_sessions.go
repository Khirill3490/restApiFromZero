package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"rest_api/internal/db/sqlc"
	"rest_api/internal/models"
	"rest_api/internal/repository"
)

type RefreshSessionStore struct {
	q *sqlc.Queries
}

func NewRefreshSessionStore(q *sqlc.Queries) *RefreshSessionStore {
	return &RefreshSessionStore{q: q}
}

func (s *RefreshSessionStore) Create(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) (*models.RefreshSession, error) {
	row, err := s.q.CreateRefreshSession(ctx, sqlc.CreateRefreshSessionParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, err
	}

	m := mapRefreshSession(row)
	return &m, nil
}

func (s *RefreshSessionStore) GetByHash(ctx context.Context, tokenHash string) (*models.RefreshSession, error) {
	row, err := s.q.GetRefreshSessionByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	m := mapRefreshSession(row)
	return &m, nil
}

func (s *RefreshSessionStore) RevokeByID(ctx context.Context, id int64) error {
	return s.q.RevokeRefreshSessionByID(ctx, id)
}

func mapRefreshSession(r sqlc.RefreshSession) models.RefreshSession {
	var revokedAt *time.Time
	if r.RevokedAt.Valid {
		t := r.RevokedAt.Time
		revokedAt = &t
	}

	return models.RefreshSession{
		ID:        r.ID,
		UserID:    r.UserID,
		TokenHash: r.TokenHash,
		ExpiresAt: r.ExpiresAt,
		CreatedAt: r.CreatedAt,
		RevokedAt: revokedAt,
	}
}
