-- name: SendMessageToGroup :one
INSERT INTO messages (sender_id, group_id, message_content, attachments, message_type, reply_to)
SELECT
    u.user_id,
    g.group_id,
    sqlc.arg(message_content)::TEXT,
    sqlc.arg(attachments)::JSONB,
    COALESCE(sqlc.arg(message_type)::INT, 1),
    sqlc.arg(reply_to)::INT
FROM users u, groups g
WHERE u.user_uuid = sqlc.arg(sender_uuid)::UUID
  AND g.group_uuid = sqlc.arg(group_uuid)::UUID
  AND EXISTS (
      SELECT 1 FROM group_members gm
      WHERE gm.group_id = g.group_id
        AND gm.member_id = u.user_id
        AND gm.left_at IS NULL
  )
  AND g.group_deleted_at IS NULL
RETURNING
    message_id,
    sender_id,
    group_id,
    message_content,
    attachments,
    message_type,
    reply_to,
    message_created_at;

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
