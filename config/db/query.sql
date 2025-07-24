-- name: GetUsers :many
SELECT * FROM users
ORDER BY (first_name || last_name);
