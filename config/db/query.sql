-- name: GetUsers :many
SELECT * FROM users
ORDER BY (first_name || last_name);

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 OR email = $1;

-- name: CreateUser :one
INSERT INTO users (
    email, 
    username, 
    password_hash, 
    password_salt,
    first_name, 
    last_name, 
    is_active, 
    role,
    password_changed_at
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9 ) RETURNING id;

-- name: LoginEvent :exec
UPDATE users 
SET last_login_at = $1
WHERE id = $2;

-- name: UserExists :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE username = $1 OR email = $2
);

