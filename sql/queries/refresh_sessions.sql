-- name: CreateRefreshSession :one
INSERT INTO refresh_sessions (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, token_hash, expires_at, created_at, revoked_at;

-- name: GetRefreshSessionByHash :one
SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
FROM refresh_sessions
WHERE token_hash = $1;

-- name: RevokeRefreshSessionByHash :exec
UPDATE refresh_sessions
SET revoked_at = now()
WHERE token_hash = $1 AND revoked_at IS NULL;
