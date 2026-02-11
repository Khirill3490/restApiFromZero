package repository

import (
	"context"
	"rest_api/internal/models"
)

type TaskRepository interface {
	ListByUser(ctx context.Context, userID int64, limit, offset int32) ([]models.Task, error)
	GetByID(ctx context.Context, userID, id int64) (*models.Task, error)
	Create(ctx context.Context, userID int64, input models.CreateTaskInput) (*models.Task, error)
	Update(ctx context.Context, userID, id int64, input models.UpdateTaskInput) (*models.Task, error)
	Delete(ctx context.Context, userID, id int64) error
}

