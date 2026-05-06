-- name: SendMessage :one
WITH sender AS(
    SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID
),
conv AS(
    SELECT conversation_id
        FROM conversations
        WHERE conversation_id = sqlc.arg(conversation_id)::INT
            AND EXISTS(
                SELECT 1 FROM conversation_participants
                WHERE conversation_id = sqlc.arg(conversation_id)::INT
                    AND user_id = (SELECT user_id FROM sender)
                    AND left_at IS NULL
            )
),
new_msg AS(
    INSERT INTO messages (conversation_id, sender_id, message_content, attachments, message_type, reply_to)
    SELECT
        c.conversation_id,s.user_id, sqlc.arg(content)::TEXT,
        COALESCE(sqlc.arg(attachments)::JSONB,'[]'::JSONB),
        sqlc.arg(message_type)::INT,
        sqlc.arg(reply_to)::INT
    FROM conv c, sender s
    WHERE EXISTS (SELECT 1 FROM conv)
    RETURNING *
)
UPDATE conversations
SET
    last_message_id = (SELECT message_id FROM new_msg),
    last_message_at = (SELECT message_created_at FROM new_msg),
    conversation_updated_at = now()
WHERE conversation_id = sqlc.arg(conversation_id)::INT
    AND EXISTS (SELECT 1 FROM new_msg)
RETURNING (SELECT message_id FROM new_msg) AS message_id;

-- name: GetMessages :many
SELECT m.message_id, m.sender_id, m.message_content, m.attachments,
       m.message_type, m.reply_to, m.message_created_at
FROM messages m
WHERE m.conversation_id = sqlc.arg(conversation_id)::INT
  AND m.message_deleted_at IS NULL
  AND (sqlc.narg(cursor_time)::timestamptz IS NULL
       OR (m.message_created_at, m.message_id) < (sqlc.narg(cursor_time)::timestamptz, sqlc.narg(cursor_id)::int))
ORDER BY m.message_created_at DESC, m.message_id DESC
LIMIT sqlc.arg(limitArg)::INT;

-- name: ValidateReply :one
SELECT EXISTS(
    SELECT 1 FROM messages
    WHERE message_id = sqlc.arg(message_id)::int
      AND conversation_id = sqlc.arg(conversation_id)::int
      AND message_deleted_at IS NULL
)::boolean;
-- name: SearchMessageInGroup :many
SELECT
    u.user_uuid,
    ms.message_content,
    ms.attachments,
    ms.message_type,
    ms.reply_to,
    ms.message_created_at
FROM messages ms
INNER JOIN groups g ON g.group_id = ms.group_id
INNER JOIN users u ON u.user_id = ms.sender_id
WHERE
    g.group_uuid = sqlc.arg(group_uuid)::UUID
    AND EXISTS (
        SELECT 1 FROM group_members gm
        WHERE gm.group_id = g.group_id
            AND gm.member_id = u.user_id
            AND gm.left_at IS NULL
    )
    AND to_tsvector('english', COALESCE(ms.message_content, '')) @@
        plainto_tsquery('english', sqlc.arg(search)::TEXT)
ORDER BY ms.message_created_at DESC
LIMIT sqlc.arg(limitArg) OFFSET sqlc.arg(offsetArg);

-- name: SoftDeleteMessage :one
UPDATE messages
SET message_deleted_at = now()
WHERE message_id = sqlc.arg(message_id)::INT
  AND sender_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID)
  AND group_id = (SELECT group_id FROM groups WHERE group_uuid = sqlc.arg(group_uuid)::UUID)
  AND EXISTS (
      SELECT 1 FROM group_members gm
      WHERE gm.group_id = messages.group_id
        AND gm.member_id = messages.sender_id
        AND gm.left_at IS NULL
  )
RETURNING *;

-- name: EditMessage :one
UPDATE messages
SET
    message_content = sqlc.arg(message_content)
WHERE
    group_id = (SELECT group_id FROM groups WHERE group_uuid = sqlc.arg(group_uuid)::UUID)
    AND sender_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID)
    AND EXISTS(
        SELECT 1
        FROM group_members gm
        INNER JOIN groups g ON g.group_id = gm.group_id
        INNER JOIN users u ON u.user_id = sender_id
        WHERE
            gm.left_at IS NULL
    )
RETURNING *;
