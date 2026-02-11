-- name: ListTasksByUser :many
SELECT id, user_id, title, description, completed, created_at, updated_at
FROM tasks
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTaskByIDByUser :one
SELECT id, user_id, title, description, completed, created_at, updated_at
FROM tasks
WHERE id = $1 AND user_id = $2;

-- name: CreateTask :one
INSERT INTO tasks (user_id, title, description, completed)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, title, description, completed, created_at, updated_at;

-- name: UpdateTaskByUser :one
UPDATE tasks
SET
  title = $3,
  description = $4,
  completed = $5,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, title, description, completed, created_at, updated_at;

-- name: DeleteTaskByUser :one
DELETE FROM tasks
WHERE id = $1 AND user_id = $2
RETURNING id;
