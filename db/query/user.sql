-- name: CreateUser :one
INSERT INTO users (
    name,
    username,
    email,
    phone,
    has_avatar,
    password,
    role_id,
    is_superuser
)
VALUES (
           $1, $2, $3, $4, $5, $6, $7, $8
       )
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users SET
                 name = $2,
                 username = $3,
                 email = $4,
                 phone = $5,
                 has_avatar = $6,
                 email_verified = $7,
                 phone_verified = $8,
                 two_factor_enabled = $9,
                 two_factor_secret = $10,
                 last_login = $11,
                 login_attempts = $12,
                 locked_until = $13,
                 password_changed_at = $14,
                 role_id = $15,
                 is_superuser = $16,
                 updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2,
    password_changed_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: RestoreUser :exec
UPDATE users
SET deleted_at = NULL, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NOT NULL;

-- name: HardDeleteUser :execrows
DELETE FROM users
WHERE id = $1;
