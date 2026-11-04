-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, issued_at, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRefreshTokenByID :one
SELECT * FROM refresh_tokens
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetRefreshTokenByIDForUpdate :one
SELECT * FROM refresh_tokens
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1
FOR UPDATE;

-- name: ListRefreshTokens :many
SELECT * FROM refresh_tokens
ORDER BY issued_at DESC
LIMIT $1 OFFSET $2;

-- name: ListRefreshTokensByUser :many
SELECT * FROM refresh_tokens
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY issued_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateRefreshToken :one
UPDATE refresh_tokens
SET
    user_id = $2,
    issued_at = $3,
    expires_at = $4
WHERE id = $1
RETURNING *;

-- name: SoftDeleteRefreshToken :exec
UPDATE refresh_tokens
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreRefreshToken :exec
UPDATE refresh_tokens
SET deleted_at = NULL
WHERE id = $1 AND deleted_at IS NOT NULL;

-- name: HardDeleteRefreshToken :execrows
DELETE FROM refresh_tokens
WHERE id = $1;

-- name: HardDeleteRefreshTokensByUser :execrows
DELETE FROM refresh_tokens
WHERE user_id = $1;