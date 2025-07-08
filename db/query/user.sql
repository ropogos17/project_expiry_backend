-- name: CreateUser :one
INSERT INTO users (
    name,
    username,
    email,
    phone,
    has_avatar,
    password,
    salt,
    role_id,
    is_superuser,
    status_id
)
VALUES (
           $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
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
                 salt = $8,
                 email_verified = $9,
                 phone_verified = $10,
                 two_factor_enabled = $11,
                 two_factor_secret = $12,
                 last_login = $13,
                 login_attempts = $14,
                 locked_until = $15,
                 password_changed_at = $16,
                 role_id = $17,
                 is_superuser = $18,
                 status_id = $19,
                 deleted_at = $20,
                 updated_at = $21
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
