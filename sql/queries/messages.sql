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

-- name: GetGroupMemberIDs :many
SELECT user_id
FROM chat_group_members
WHERE chat_group_id = $1;

-- name: GetConversationsByUserID :many
WITH recent AS (
    SELECT
        CASE
            WHEN groups.id IS NOT NULL THEN messages.receiver_id
            WHEN messages.sender_id = $1 THEN messages.receiver_id
            ELSE messages.sender_id
        END AS other_id,
        MAX(messages.created_at) AS latest
    FROM messages
    LEFT JOIN chat_groups AS groups ON groups.id = messages.receiver_id
    WHERE (
        groups.id IS NOT NULL
        AND messages.receiver_id IN (
            SELECT chat_group_id
            FROM chat_group_members
            WHERE user_id = $1
        )
    ) OR (
        groups.id IS NULL
        AND (messages.sender_id = $1 OR messages.receiver_id = $1)
    )
    GROUP BY other_id
)
SELECT u.id, u.name, r.latest
FROM recent r
JOIN users u ON r.other_id = u.id

UNION ALL

SELECT c.id, c.name, r.latest
FROM chat_group_members AS members
JOIN chat_groups AS c ON c.id = members.chat_group_id
LEFT JOIN recent AS r ON r.other_id = c.id
WHERE members.user_id = $1

ORDER BY latest DESC NULLS LAST;

-- name: GetGroupMessages :many
SELECT *
FROM messages 
WHERE receiver_id = $1
ORDER BY created_at DESC;


-- name: GetDuelMessages :many

SELECT * 
FROM messages 
WHERE (receiver_id = $1 AND sender_id = $2) OR (receiver_id = $2 AND sender_id = $1)
ORDER BY created_at DESC;

