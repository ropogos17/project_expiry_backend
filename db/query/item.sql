-- name: CreateItem :one
INSERT INTO items (
    name, category_id, description, quantity, has_image, expires_at, storage_id, shelf, last_modified
)
VALUES (
           $1, $2, $3, $4, $5, $6,
           $7, $8, $9
       )
RETURNING *;

-- name: GetItem :one
SELECT * FROM items
WHERE id = $1 LIMIT 1;

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
    is_deleted = $7,
    stored_at = $8,
    updated_at = $9,
    expires_at = $10,
    storage_id = $11,
    shelf = $12,
    last_modified = $13
WHERE id = $1
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = $1;
