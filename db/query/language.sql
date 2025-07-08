-- name: CreateLanguage :one
INSERT INTO languages (
    en, hu, ro, es, zh, fr, de, pt, ru, ua, ar
)
VALUES (
           $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
       )
RETURNING *;

-- name: GetLanguage :one
SELECT * FROM languages
WHERE id = $1 LIMIT 1;

-- name: ListLanguages :many
SELECT * FROM languages
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateLanguage :one
UPDATE languages
SET
    en = $2,
    hu = $3,
    ro = $4,
    es = $5,
    zh = $6,
    fr = $7,
    de = $8,
    pt = $9,
    ru = $10,
    ua = $11,
    ar = $12
WHERE id = $1
RETURNING *;

-- name: DeleteLanguage :exec
DELETE FROM languages
WHERE id = $1;
