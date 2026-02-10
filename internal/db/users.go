package db

import (
	"context"
	"database/sql"
	"errors"

	"rest_api/internal/db/sqlc"
	"rest_api/internal/models"
	"rest_api/internal/repository"
)

type UserStore struct {
	q *sqlc.Queries
}

func NewUserStore(q *sqlc.Queries) *UserStore {
	return &UserStore{q: q}
}

func (s *UserStore) Create(ctx context.Context, email, passwordHash string) (*models.User, error) {
	u, err := s.q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, repository.ErrConflict
		}
		return nil, err
	}

	return mapUser(u), nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	u, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return mapUser(u), nil
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*models.User, error) {
	u, err := s.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return mapUser(u), nil
}

func mapUser(u sqlc.User) *models.User {
	return &models.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
	}
}
