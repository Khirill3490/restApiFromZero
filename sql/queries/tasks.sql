-- name: ListTasks :many
SELECT id, title, description, completed, created_at, updated_at
FROM tasks
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTaskByID :one
SELECT id, title, description, completed, created_at, updated_at
FROM tasks
WHERE id = $1;

-- name: CreateTask :one
INSERT INTO tasks (title, description, completed)
VALUES ($1, $2, $3)
RETURNING id, title, description, completed, created_at, updated_at;

-- name: UpdateTask :one
UPDATE tasks
SET
  title = $2,
  description = $3,
  completed = $4,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, title, description, completed, created_at, updated_at;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;
