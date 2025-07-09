-- name: CreateUserStatus :one
INSERT INTO user_statuses (name)
VALUES ($1)
RETURNING *;

-- name: GetUserStatus :one
SELECT * FROM user_statuses
WHERE id = $1 LIMIT 1;

-- name: ListUserStatuses :many
SELECT * FROM user_statuses
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUserStatus :one
UPDATE user_statuses
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUserStatus :exec
DELETE FROM user_statuses
WHERE id = $1;