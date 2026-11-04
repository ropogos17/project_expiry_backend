-- name: CreateLog :one
INSERT INTO logs (action_type_id, message, data, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetLogByID :one
SELECT * FROM logs
WHERE id = $1 LIMIT 1;

-- name: ListLogs :many
SELECT * FROM logs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateLog :one
UPDATE logs
SET
    action_type_id = $2,
    message = $3,
    data = $4
WHERE id = $1
RETURNING *;

-- name: HardDeleteLog :execrows
DELETE FROM logs
WHERE id = $1;
