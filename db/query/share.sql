-- name: CreateShare :one
INSERT INTO shares (
    source_user_id, storage_id, target_user_id, can_share
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetShareByID :one
SELECT * FROM shares
WHERE id = $1 LIMIT 1;

-- name: GetShareByStorageAndTargetUser :one
SELECT * FROM shares
WHERE storage_id = $1 AND target_user_id = $2 LIMIT 1;

-- name: ListShares :many
SELECT * FROM shares
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListSharesByUser :many
SELECT * FROM shares
WHERE source_user_id = $3
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateShare :one
UPDATE shares
SET
    source_user_id = $2,
    storage_id = $3,
    target_user_id = $4,
    can_share = $5,
    shared_at = $6
WHERE id = $1
RETURNING *;

-- name: HardDeleteShare :execrows
DELETE FROM shares
WHERE id = $1;
