-- name: CreateStatus :one
INSERT INTO user_statuses (name)
VALUES ($1)
RETURNING *;

-- name: GetStatus :one
SELECT * FROM user_statuses
WHERE id = $1 LIMIT 1;

-- name: ListStatuses :many
SELECT * FROM user_statuses
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateStatus :one
UPDATE user_statuses
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteStatus :exec
DELETE FROM user_statuses
WHERE id = $1;