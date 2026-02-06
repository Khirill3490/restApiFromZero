package db

import (
	"context"
	"database/sql"
	"rest_api/internal/models"

	"rest_api/internal/db/sqlc"
)

type TaskStore struct {
	db *sql.DB
	q  *sqlc.Queries
}

func NewTaskStore(db *sql.DB) *TaskStore {
	return &TaskStore{
		db: db,
		q:  sqlc.New(db),
	}
}

func (ts *TaskStore) List(ctx context.Context, limit, offset int32) ([]models.Task, error) {
	rows, err := ts.q.ListTasks(ctx, sqlc.ListTasksParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	tasks := make([]models.Task, 0, len(rows))

	for _, r := range rows {
		tasks = append(tasks, models.Task{
			ID:          r.ID,
			Title:       r.Title,
			Description: r.Description,
			Completed:   r.Completed,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		})
	}

	return tasks, nil
}

func (ts *TaskStore) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	r, err := ts.q.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task := mapToModelTask(r)

	return task, nil
}

func (ts *TaskStore) Create(ctx context.Context, input models.CreateTaskInput) (*models.Task, error) {
	r, err := ts.q.CreateTask(ctx, sqlc.CreateTaskParams{
		Title:       input.Title,
		Description: input.Description,
		Completed:   input.Completed,
	})
	if err != nil {
		return nil, err
	}

	task := mapToModelTask(r)

	return task, nil
}

func (ts *TaskStore) Update(ctx context.Context, id int64, input models.UpdateTaskInput) (*models.Task, error) {
	if input.Title == nil || input.Description == nil || input.Completed == nil {
		return nil, sql.ErrNoRows
	}

	cur, err := ts.q.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	title := cur.Title
	if input.Title != nil {
		title = *input.Title
	}

	desc := cur.Description
	if input.Description != nil {
		desc = *input.Description
	}

	completed := cur.Completed
	if input.Completed != nil {
		completed = *input.Completed
	}

	r, err := ts.q.UpdateTask(ctx, sqlc.UpdateTaskParams{
		ID:          id,
		Title:       title,
		Description: desc,
		Completed:   completed,
	})
	if err != nil {
		return nil, err
	}

	task := mapToModelTask(r)

	return task, nil
}

func (ts *TaskStore) Delete(ctx context.Context, id int64) error {
	return ts.q.DeleteTask(ctx, id)
}

func mapToModelTask(r sqlc.Task) *models.Task {
	return &models.Task{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		Completed:   r.Completed,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}
