-- name: CreateActionType :one
INSERT INTO action_types (name)
VALUES ($1)
RETURNING *;

-- name: GetActionType :one
SELECT * FROM action_types
WHERE id = $1 LIMIT 1;

-- name: ListActionTypes :many
SELECT * FROM action_types
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateActionType :one
UPDATE action_types
SET name = $2, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteActionType :one
UPDATE action_types
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RestoreActionType :one
UPDATE action_types
SET deleted_at = NULL, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NOT NULL
RETURNING *;

-- name: HardDeleteActionType :exec
DELETE FROM action_types
WHERE id = $1;