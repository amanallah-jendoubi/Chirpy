-- name: GetChirpById :one
SELECT
    id,
    created_at,
    updated_at,
    user_id,
    body
FROM chirps
WHERE id = $1 ;