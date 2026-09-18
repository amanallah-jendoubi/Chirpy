-- name: CreateChirps :one
INSERT INTO chirps (created_at, updated_at, body, user_id)
VALUES (
    now(), now(), $1, $2
)
RETURNING *;