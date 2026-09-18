-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, name)
VALUES (
    now(), now(), $1, $2
)
RETURNING *;