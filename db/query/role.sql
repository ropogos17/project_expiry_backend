-- name: CreateRole :one
INSERT INTO roles (name)
VALUES ($1)
RETURNING *;

-- name: GetRoleByID :one
SELECT * FROM roles
WHERE id = $1 LIMIT 1;

-- name: RoleExistsByName :one
SELECT EXISTS (
    SELECT 1 FROM roles
    WHERE name = $1
);

-- name: ListRoles :many
SELECT * FROM roles
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateRole :one
UPDATE roles
SET name = $2
WHERE id = $1
RETURNING *;

-- name: HardDeleteRole :execrows
DELETE FROM roles
WHERE id = $1;