-- name: CreateCategory :one
INSERT INTO categories (name)
VALUES ($1)
RETURNING *;

-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: CategoryExistsByName :one
SELECT EXISTS (
    SELECT 1 FROM categories
    WHERE name = $1
);

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2
WHERE id = $1
RETURNING *;

-- name: HardDeleteCategory :execrows
DELETE FROM categories
WHERE id = $1;