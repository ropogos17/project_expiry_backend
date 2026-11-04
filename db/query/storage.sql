-- name: CreateStorage :one
INSERT INTO storages (
    name, description, location, shelf_count, has_image, owner
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetStorageByID :one
SELECT * FROM storages
WHERE id = $1 LIMIT 1;

-- name: ListStorages :many
SELECT * FROM storages
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListStoragesByOwner :many
SELECT * FROM storages
WHERE owner = $3 AND deleted_at IS NULL
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateStorage :one
UPDATE storages
SET
    name = $2,
    description = $3,
    location = $4,
    shelf_count = $5,
    has_image = $6,
    owner = $7,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteStorage :exec
UPDATE storages
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreStorage :exec
UPDATE storages
SET deleted_at = NULL, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NOT NULL;

-- name: HardDeleteStorage :execrows
DELETE FROM storages
WHERE id = $1;
