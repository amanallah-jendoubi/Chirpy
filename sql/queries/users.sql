-- name: CreateUser :one
INSERT INTO users (password, name)
VALUES ($1, $2)
RETURNING *;

-- name: UserExists :one
SELECT EXISTS(
    SELECT 1 FROM users WHERE name = $1
) AS exists;

-- name: GetUserByName :one
SELECT *
FROM users
WHERE name = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT name
FROM users
WHERE id = $1;