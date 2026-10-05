-- name: CreateChatGroupWithCreator :one
WITH new_group AS (
    INSERT INTO chat_groups (name, created_by)
    VALUES ($1, $2)
    RETURNING id, name, created_at, created_by
), added_creator AS (
    INSERT INTO chat_group_members (user_id, chat_group_id)
    SELECT created_by, id
    FROM new_group
    RETURNING chat_group_id
)
SELECT id, name, created_at, created_by
FROM new_group;

-- name: AddChatGroupMember :exec
INSERT INTO chat_group_members (user_id, chat_group_id)
VALUES ($1, $2)
ON CONFLICT (user_id, chat_group_id) DO NOTHING;