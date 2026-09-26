-- name: CreateMessage :one
INSERT INTO messages (
    sender_id,
    receiver_id,
    body
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: ChatGroupExists :one
SELECT EXISTS (
    SELECT 1 FROM chat_groups WHERE id = $1
) AS exists;

-- name: IsGroupMember :one
SELECT EXISTS (
    SELECT 1
    FROM chat_group_members
    WHERE chat_group_id = $1
      AND user_id = $2
) AS exists;

-- name: GetConversationsByUserID :many
WITH recent AS (
    SELECT receiver_id, MAX(created_at) AS latest
    FROM messages
    WHERE sender_id = $1
    GROUP BY receiver_id
)
SELECT u.name, r.latest
FROM recent r , users u 
WHERE r.receiver_id = u.id 

UNION ALL

SELECT c.name, r.latest
FROM recent r, chat_groups c
WHERE r.receiver_id = c.id 

ORDER BY latest DESC;

