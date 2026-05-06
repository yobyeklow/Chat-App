-- name: FindDMConversation :one
SELECT c.conversation_id
FROM conversations c
WHERE c.conversation_type = 1
  AND EXISTS (
    SELECT 1 FROM conversation_participants cp
    WHERE cp.conversation_id = c.conversation_id AND cp.user_id = sqlc.arg(user1_id)
  )
  AND EXISTS (
    SELECT 1 FROM conversation_participants cp
    WHERE cp.conversation_id = c.conversation_id AND cp.user_id = sqlc.arg(user2_id)
  )
LIMIT 1;

-- name: CreateDMConversation :one
INSERT INTO conversations (conversation_type, reference_id)
VALUES (1, nextval('dm_reference_seq'))
RETURNING conversation_id;

-- name: AddParticipantToConversation :exec
INSERT INTO conversation_participants (conversation_id, user_id)
VALUES (sqlc.arg(conversation_id), sqlc.arg(user_id))
ON CONFLICT (conversation_id, user_id) WHERE left_at IS NULL DO NOTHING;

-- name: HardDeleteConversation :exec
DELETE FROM conversations
WHERE conversation_type = 2
    AND reference_id = (SELECT group_id FROM groups WHERE group_uuid = sqlc.arg(group_uuid));

-- name: LeaveConversation :exec
UPDATE conversation_participants
SET left_at = now()
WHERE conversation_id = sqlc.arg(conversation_id)::int
    AND user_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID)
    AND left_at IS NULL;

-- name: GetDMConversation :many
SELECT
    c.conversation_id,
    1 AS conversation_type,
    u.user_email,
    c.last_message_id,
    c.last_message_at,
    (
        SELECT COUNT(*)
        FROM messages m
        WHERE m.conversation_id = c.conversation_id
            AND m.message_id > COALESCE(
                (SELECT last_read_message_id FROM conversation_participants
                    WHERE conversation_id = c.conversation_id AND user_id = sqlc.arg(cur_user_id)::INT)
            ,0)
            AND m.message_deleted_at IS NULL
    ) AS unread_count
FROM conversations c
JOIN conversation_participants cp ON c.conversation_id = cp.conversation_id AND cp.user_id <> sqlc.arg(cur_user_id)::INT
JOIN users u ON cp.user_id = u.user_id
WHERE c.conversation_type = 1
    AND EXISTS
        (
            SELECT 1 FROM conversation_participants
            WHERE conversation_id = c.conversation_id
                AND user_id = sqlc.arg(cur_user_id)::INT
                AND left_at IS NULL
        )
ORDER BY c.last_message_at DESC NULLS LAST;
-- name: GetGroupConversation :many
SELECT
    c.conversation_id,
    2 AS conversation_type,
    g.group_name AS title,
    c.last_message_id,
    c.last_message_at,
    (
        SELECT COUNT(*)
        FROM messages m
        WHERE m.conversation_id = c.conversation_id
            AND m.message_id > COALESCE(
                (SELECT last_read_message_id FROM conversation_participants
                    WHERE conversation_id = c.conversation_id AND user_id = sqlc.arg(cur_user_id)::INT)
            ,0)
            AND m.message_deleted_at IS NULL
    ) AS unread_count
FROM conversations c
JOIN groups g ON c.reference_id = g.group_id
WHERE c.conversation_type = 2
    AND EXISTS (
        SELECT 1 FROM conversation_participants
        WHERE conversation_id = c.conversation_id
            AND user_id = sqlc.arg(cur_user_id)
            AND left_at IS NULL
    )
ORDER BY c.last_message_at DESC NULLS LAST;
