-- name: GetUsers :many
SELECT * FROM users
ORDER BY (first_name || last_name);

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 OR email = $1;
