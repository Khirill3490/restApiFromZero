package repository

import (
	"context"
	"rest_api/internal/models"
)

type TaskRepository interface {
	List(ctx context.Context, limit, offset int32) ([]models.Task, error)
	GetByID(ctx context.Context, id int64) (*models.Task, error)
	Create(ctx context.Context, input models.CreateTaskInput) (*models.Task, error)
	Update(ctx context.Context, id int64, input models.UpdateTaskInput) (*models.Task, error)
	Delete(ctx context.Context, id int64) error
}
