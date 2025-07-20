-- name: CreateStorage :one
INSERT INTO storages (
    name, description, location, shelves, has_image, owner
)
VALUES (
           $1, $2, $3, $4, $5, $6
       )
RETURNING *;

-- name: GetStorage :one
SELECT * FROM storages
WHERE id = $1 LIMIT 1;

-- name: ListStorages :many
SELECT * FROM storages
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListStoragesByOwner :many
SELECT * FROM storages
WHERE owner = $3 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateStorage :one
UPDATE storages
SET
    name = $2,
    description = $3,
    location = $4,
    shelves = $5,
    has_image = $6,
    owner = $7,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteStorage :one
UPDATE storages
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RestoreStorage :one
UPDATE storages
SET deleted_at = NULL, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NOT NULL
RETURNING *;

-- name: HardDeleteStorage :exec
DELETE FROM storages
WHERE id = $1;
