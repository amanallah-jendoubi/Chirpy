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


-- name: GetRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1
LIMIT 1;
-- name: RevokeRefreshTokensByFamilyID :exec
DELETE FROM refresh_tokens
WHERE family_id = $1;