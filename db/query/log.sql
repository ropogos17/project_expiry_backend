-- name: CreateLog :one
INSERT INTO logs (
    action_type, message, data, user_id
)
VALUES (
           $1, $2, $3, $4
       )
RETURNING *;

-- name: GetLog :one
SELECT * FROM logs
WHERE id = $1 LIMIT 1;

-- name: ListLogs :many
SELECT * FROM logs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateLog :one
UPDATE logs
SET
    created_at = $2,
    action_type = $3,
    message = $4,
    data = $5,
    user_id = $6
WHERE id = $1
RETURNING *;

-- name: DeleteLog :exec
DELETE FROM logs
WHERE id = $1;
