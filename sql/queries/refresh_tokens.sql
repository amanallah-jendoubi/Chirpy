-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id,
    token,
    expires_at,
    family_id
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;