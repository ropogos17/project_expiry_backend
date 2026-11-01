-- name: CreateAnnouncement :one
INSERT INTO announcements (name, description, external_link, created_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAnnouncementByID :one
SELECT * FROM announcements
WHERE id = $1 LIMIT 1;

-- name: ListAnnouncements :many
SELECT * FROM announcements
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateAnnouncement :one
UPDATE announcements
SET
    name = $2,
    description = $3,
    external_link = $4,
    created_by = $5,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteAnnouncement :exec
UPDATE announcements
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreAnnouncement :exec
UPDATE announcements
SET deleted_at = NULL, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NOT NULL;

-- name: HardDeleteAnnouncement :exec
DELETE FROM announcements
WHERE id = $1;