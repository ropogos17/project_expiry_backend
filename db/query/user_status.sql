-- name: CreateUserStatus :one
INSERT INTO user_statuses (name)
VALUES ($1)
RETURNING *;

-- name: GetUserStatus :one
SELECT * FROM user_statuses
WHERE id = $1 LIMIT 1;

-- name: ListUserStatuses :many
SELECT * FROM user_statuses
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUserStatus :one
UPDATE user_statuses
SET name = $2, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteUserStatus :one
UPDATE user_statuses
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RestoreUserStatus :one
UPDATE user_statuses
SET deleted_at = NULL, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NOT NULL
RETURNING *;

-- name: HardDeleteUserStatus :exec
DELETE FROM user_statuses
WHERE id = $1;