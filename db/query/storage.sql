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
WHERE owner = $3
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: GetShareByStorageAndOwner :one
SELECT * FROM shares
WHERE storage_id = $1 AND target_user_id = $2 LIMIT 1;

-- name: UpdateStorage :one
UPDATE storages
SET
    name = $2,
    description = $3,
    location = $4,
    shelves = $5,
    has_image = $6,
    owner = $7
WHERE id = $1
RETURNING *;

-- name: DeleteStorage :exec
UPDATE storages
SET is_deleted = true, updated_at = NOW()
WHERE id = $1;

-- name: DeleteStorage :exec
DELETE FROM storages
WHERE id = $1;
