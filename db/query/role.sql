-- name: CreateRole :one
INSERT INTO roles (
    name, bio
) VALUES (
             $1, $2
         )
RETURNING *;