-- name: CreateShoppingListItem :one
INSERT INTO shopping_list_items (
    name, description, quantity, is_checked, is_deleted,
    updated_at, last_modified
)
VALUES (
           $1, $2, $3, $4, $5,
           $6, $7
       )
RETURNING *;

-- name: GetShoppingListItem :one
SELECT * FROM shopping_list_items
WHERE id = $1 LIMIT 1;

-- name: ListShoppingListItems :many
SELECT * FROM shopping_list_items
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateShoppingListItem :one
UPDATE shopping_list_items
SET
    name = $2,
    description = $3,
    quantity = $4,
    is_checked = $5,
    is_deleted = $6,
    created_at = $7,
    updated_at = $8,
    last_modified = $9
WHERE id = $1
RETURNING *;

-- name: DeleteShoppingListItem :exec
DELETE FROM shopping_list_items
WHERE id = $1;
