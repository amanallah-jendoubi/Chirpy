-- name: CreateUser :one
INSERT INTO users (password, name)
VALUES ($1, $2)
RETURNING *;

-- name: UserExists :one
SELECT EXISTS(
    SELECT 1 FROM users WHERE name = $1
) AS exists;