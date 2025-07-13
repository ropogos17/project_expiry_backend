-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id, issued_at, expires_at
)
VALUES (
           $1, $2, $3
       )
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE id = $1
LIMIT 1;

-- name: ListRefreshTokens :many
SELECT * FROM refresh_tokens
ORDER BY issued_at DESC
LIMIT $1 OFFSET $2;

-- name: ListRefreshTokensByUser :many
SELECT * FROM refresh_tokens
WHERE user_id = $1
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

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE id = $1;
