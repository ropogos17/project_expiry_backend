-- name: CreateAnnouncement :one
INSERT INTO announcements (
    name, description, external_link, has_read
)
VALUES (
           $1, $2, $3, $4
       )
RETURNING *;

-- name: GetAnnouncement :one
SELECT * FROM announcements
WHERE id = $1 LIMIT 1;

-- name: ListAnnouncements :many
SELECT * FROM announcements
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateAnnouncement :one
UPDATE announcements
SET
    name = $2,
    description = $3,
    external_link = $4,
    has_read = $5
WHERE id = $1
RETURNING *;

-- name: DeleteAnnouncement :exec
DELETE FROM announcements
WHERE id = $1;
