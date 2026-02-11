package db

import (
	"context"
	"database/sql"
	"errors"

	"rest_api/internal/db/sqlc"
	"rest_api/internal/models"
	"rest_api/internal/repository"
)

type TaskStore struct {
	q *sqlc.Queries
}

func NewTaskStore(q *sqlc.Queries) *TaskStore {
	return &TaskStore{q: q}
}

// ListByUser возвращает задачи только конкретного пользователя
func (ts *TaskStore) ListByUser(ctx context.Context, userID int64, limit, offset int32) ([]models.Task, error) {
	rows, err := ts.q.ListTasksByUser(ctx, sqlc.ListTasksByUserParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	tasks := make([]models.Task, 0, len(rows))
	for _, r := range rows {
		tasks = append(tasks, mapToModelTask(r))
	}

	return tasks, nil
}

// GetByID возвращает задачу только если она принадлежит этому userID
func (ts *TaskStore) GetByID(ctx context.Context, userID, id int64) (*models.Task, error) {
	r, err := ts.q.GetTaskByIDByUser(ctx, sqlc.GetTaskByIDByUserParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	task := mapToModelTask(r)
	return &task, nil
}

// Create создаёт задачу и сразу привязывает к userID
func (ts *TaskStore) Create(ctx context.Context, userID int64, input models.CreateTaskInput) (*models.Task, error) {
	r, err := ts.q.CreateTask(ctx, sqlc.CreateTaskParams{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		Completed:   input.Completed,
	})
	if err != nil {
		return nil, err
	}

	task := mapToModelTask(r)
	return &task, nil
}

// Update обновляет задачу только если она принадлежит userID.
// Логика частичного обновления остаётся как у тебя: сначала читаем текущую, потом мерджим поля.
func (ts *TaskStore) Update(ctx context.Context, userID, id int64, input models.UpdateTaskInput) (*models.Task, error) {
	if input.Title == nil && input.Description == nil && input.Completed == nil {
		return nil, errors.New("no fields to update")
	}

	cur, err := ts.q.GetTaskByIDByUser(ctx, sqlc.GetTaskByIDByUserParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
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

	r, err := ts.q.UpdateTaskByUser(ctx, sqlc.UpdateTaskByUserParams{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: desc,
		Completed:   completed,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	task := mapToModelTask(r)
	return &task, nil
}

// Delete удаляет задачу только если она принадлежит userID
func (ts *TaskStore) Delete(ctx context.Context, userID, id int64) error {
	_, err := ts.q.DeleteTaskByUser(ctx, sqlc.DeleteTaskByUserParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}

// ВАЖНО:
// этот маппер зависит от того, какой тип sqlc сгенерировал для UserID.
//
// Вариант А (чаще всего, если user_id NOT NULL): r.UserID — int64
// Вариант Б (если user_id nullable): r.UserID — sql.NullInt64 -> тогда нужно брать r.UserID.Int64
func mapToModelTask(r any) models.Task {
	switch t := r.(type) {

	// --- Вариант А: sqlc.Task, где UserID int64
	case sqlc.Task:
		return models.Task{
			ID:          t.ID,
			UserID:      t.UserID,
			Title:       t.Title,
			Description: t.Description,
			Completed:   t.Completed,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}

	// --- Вариант Б: sqlc.Task, где UserID sql.NullInt64
	// (если у тебя так — раскомментируй этот кейс и поправь модель Task под UserID)
	/*
	case sqlc.Task:
		var uid int64
		if t.UserID.Valid {
			uid = t.UserID.Int64
		}
		return models.Task{
			ID:          t.ID,
			UserID:      uid,
			Title:       t.Title,
			Description: t.Description,
			Completed:   t.Completed,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
	*/

	default:
		panic("unexpected task row type in mapToModelTask")
	}
}
