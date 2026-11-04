-- name: CreateItem :one
INSERT INTO items (
    name, category_id, description, quantity, has_image, stored_at, expires_at, storage_id, shelf, last_modified_by
)
VALUES ($1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10)
RETURNING *;

-- name: GetItemByID :one
SELECT * FROM items
WHERE id = $1 LIMIT 1;

-- name: ListItemsByStorage :many
SELECT * FROM items
WHERE storage_id = $1 AND deleted_at IS NULL
ORDER BY expires_at
LIMIT $2 OFFSET $3;

-- name: ListItemsByUser :many
SELECT DISTINCT i.* FROM items i
JOIN storages s ON i.storage_id = s.id
LEFT JOIN shares sh ON s.id = sh.storage_id
WHERE (s.owner = $1 OR sh.target_user_id = $1) AND i.deleted_at IS NULL
ORDER BY i.expires_at
LIMIT $2 OFFSET $3;

-- name: ListActiveItemsByStorage :many
SELECT * FROM items
WHERE storage_id = $1
  AND deleted_at IS NULL
  AND expires_at > $2
ORDER BY expires_at
LIMIT $3 OFFSET $4;

-- name: ListActiveItemsByUser :many
SELECT DISTINCT i.* FROM items i
JOIN storages s ON i.storage_id = s.id
LEFT JOIN shares sh ON s.id = sh.storage_id
WHERE (s.owner = $1 OR sh.target_user_id = $1)
  AND i.deleted_at IS NULL
  AND i.expires_at > $2
ORDER BY i.expires_at
LIMIT $3 OFFSET $4;

-- name: ListExpiredItemsByStorage :many
SELECT * FROM items
WHERE storage_id = $1
  AND deleted_at IS NULL
  AND expires_at <= $2
ORDER BY expires_at
LIMIT $3 OFFSET $4;

-- name: ListExpiredItemsByUser :many
SELECT DISTINCT i.* FROM items i
JOIN storages s ON i.storage_id = s.id
LEFT JOIN shares sh ON s.id = sh.storage_id
WHERE (s.owner = $1 OR sh.target_user_id = $1)
  AND i.deleted_at IS NULL
  AND i.expires_at <= $2
ORDER BY i.expires_at
LIMIT $3 OFFSET $4;

-- name: ListDeletedItemsByStorage :many
SELECT * FROM items
WHERE storage_id = $1
  AND deleted_at IS NOT NULL
ORDER BY expires_at
LIMIT $2 OFFSET $3;

-- name: ListDeletedItemsByUser :many
SELECT DISTINCT i.* FROM items i
JOIN storages s ON i.storage_id = s.id
LEFT JOIN shares sh ON s.id = sh.storage_id
WHERE (s.owner = $1 OR sh.target_user_id = $1)
  AND i.deleted_at IS NOT NULL
ORDER BY i.expires_at
LIMIT $2 OFFSET $3;

-- name: ListItems :many
SELECT * FROM items
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateItem :one
UPDATE items
SET
    name = $2,
    category_id = $3,
    description = $4,
    quantity = $5,
    has_image = $6,
    stored_at = $7,
    updated_at = NOW(),
    expires_at = $8,
    storage_id = $9,
    shelf = $10,
    last_modified_by = $11
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteItem :exec
UPDATE items
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreItem :exec
UPDATE items
SET deleted_at = NULL, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NOT NULL;

-- name: HardDeleteItem :execrows
DELETE FROM items
WHERE id = $1;
