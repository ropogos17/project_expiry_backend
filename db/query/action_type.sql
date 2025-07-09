-- name: CreateActionType :one
INSERT INTO action_types (name)
VALUES ($1)
RETURNING *;

-- name: GetActionType :one
SELECT * FROM action_types
WHERE id = $1 LIMIT 1;

-- name: ListActionTypes :many
SELECT * FROM action_types
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateActionType :one
UPDATE action_types
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteActionType :exec
DELETE FROM action_types
WHERE id = $1;