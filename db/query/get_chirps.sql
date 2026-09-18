-- name: GetAllChirps :many
SELECT
    id,
    created_at,
    updated_at,
    user_id,
    body
FROM chirps
ORDER BY created_at ASC;
