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


// package db

// import (
// 	"context"
// 	"database/sql"
// 	"errors"

// 	"rest_api/internal/db/sqlc"
// 	"rest_api/internal/models"
// 	"rest_api/internal/repository"
// )

// type TaskStore struct {
// 	db *sql.DB
// 	q  *sqlc.Queries
// }

// func NewTaskStore(db *sql.DB) *TaskStore {
// 	return &TaskStore{
// 		db: db,
// 		q:  sqlc.New(db),
// 	}
// }

// // compile-time check: TaskStore реализует интерфейс
// var _ repository.TaskRepository = (*TaskStore)(nil)

// func (ts *TaskStore) List(ctx context.Context, limit, offset int32) ([]models.Task, error) {
// 	rows, err := ts.q.ListTasks(ctx, sqlc.ListTasksParams{
// 		Limit:  limit,
// 		Offset: offset,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	tasks := make([]models.Task, 0, len(rows))
// 	for _, r := range rows {
// 		tasks = append(tasks, models.Task{
// 			ID:          r.ID,
// 			Title:       r.Title,
// 			Description: r.Description,
// 			Completed:   r.Completed,
// 			CreatedAt:   r.CreatedAt,
// 			UpdatedAt:   r.UpdatedAt,
// 		})
// 	}
// 	return tasks, nil
// }

// func (ts *TaskStore) GetByID(ctx context.Context, id int64) (*models.Task, error) {
// 	r, err := ts.q.GetTaskByID(ctx, id)
// 	if err != nil {
// 		// сюда можно позже добавить repository.ErrNotFound, если хочешь
// 		return nil, err
// 	}

// 	task := &models.Task{
// 		ID:          r.ID,
// 		Title:       r.Title,
// 		Description: r.Description,
// 		Completed:   r.Completed,
// 		CreatedAt:   r.CreatedAt,
// 		UpdatedAt:   r.UpdatedAt,
// 	}
// 	return task, nil
// }

// func (ts *TaskStore) Create(ctx context.Context, input models.CreateTaskInput) (*models.Task, error) {
// 	r, err := ts.q.CreateTask(ctx, sqlc.CreateTaskParams{
// 		Title:       input.Title,
// 		Description: input.Description,
// 		Completed:   input.Completed,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	task := &models.Task{
// 		ID:          r.ID,
// 		Title:       r.Title,
// 		Description: r.Description,
// 		Completed:   r.Completed,
// 		CreatedAt:   r.CreatedAt,
// 		UpdatedAt:   r.UpdatedAt,
// 	}
// 	return task, nil
// }

// func (ts *TaskStore) Update(ctx context.Context, id int64, input models.UpdateTaskInput) (*models.Task, error) {
// 	// так как у тебя указатели в UpdateTaskInput, надо решить: patch или требуем все поля
// 	// Для простоты: требуем все поля (иначе нужна отдельная Patch-логика).
// 	if input.Title == nil || input.Description == nil || input.Completed == nil {
// 		return nil, errors.New("title, description and completed are required for update")
// 	}

// 	r, err := ts.q.UpdateTask(ctx, sqlc.UpdateTaskParams{
// 		ID:          id,
// 		Title:       *input.Title,
// 		Description: *input.Description,
// 		Completed:   *input.Completed,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	task := &models.Task{
// 		ID:          r.ID,
// 		Title:       r.Title,
// 		Description: r.Description,
// 		Completed:   r.Completed,
// 		CreatedAt:   r.CreatedAt,
// 		UpdatedAt:   r.UpdatedAt,
// 	}
// 	return task, nil
// }

// func (ts *TaskStore) Delete(ctx context.Context, id int64) error {
// 	return ts.q.DeleteTask(ctx, id)
// }

