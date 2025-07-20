-- name: CreateUser :one
INSERT INTO users (
    name,
    username,
    email,
    phone,
    has_avatar,
    password,
    role_id,
    is_superuser,
    status_id
)
VALUES (
           $1, $2, $3, $4, $5, $6, $7, $8, $9
       )
RETURNING *;

-- name: GetUser :one
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
                 password = $7,
                 email_verified = $8,
                 phone_verified = $9,
                 two_factor_enabled = $10,
                 two_factor_secret = $11,
                 last_login = $12,
                 login_attempts = $13,
                 locked_until = $14,
                 password_changed_at = $15,
                 role_id = $16,
                 is_superuser = $17,
                 status_id = $18,
                 deleted_at = $19,
                 updated_at = $20
WHERE id = $1 AND deleted_at IS NOT NULL
RETURNING *;

-- name: DeleteUser :one
UPDATE users
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RestoreUser :one
UPDATE users
SET deleted_at = NULL, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NOT NULL
RETURNING *;

-- name: HardDeleteUser :exec
DELETE FROM users
WHERE id = $1;
