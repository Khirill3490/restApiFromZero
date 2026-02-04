package db

import (
	"context"
	"database/sql"
	"rest_api/internal/models"
)

type TaskStore struct {
	db *sql.DB
}

func NewTaskStore(db *sql.DB) *TaskStore {
	return &TaskStore{db: db}
}

func (ts *TaskStore) List(ctx context.Context, limit, offset int32) ([]models.Task, error) {
	var tasks []models.Task
	return tasks, nil
}

func (ts *TaskStore) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	var task models.Task
	return &task, nil
}

func (ts *TaskStore) Create(ctx context.Context, input models.CreateTaskInput) (*models.Task, error) {
	var task models.Task
	return &task, nil
}

func (ts *TaskStore) Update(ctx context.Context, id int64, input models.UpdateTaskInput) (*models.Task, error) {
	var task models.Task
	return &task, nil
}

func (ts *TaskStore) Delete(ctx context.Context, id int64) error {
	return nil
}
