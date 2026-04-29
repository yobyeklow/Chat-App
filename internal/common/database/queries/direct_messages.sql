-- name: SentDirectMessages :one
INSERT INTO direct_messages(sender_id, receiver_id, dm_content, dm_attachments, dm_message_type)
SELECT
    u.user_id,
    ou.user_id,
    sqlc.arg(dm_content)::TEXT,
    sqlc.arg(dm_attachments)::JSONB,
    COALESCE(sqlc.arg(dm_message_type)::INT, 1)
FROM users u, users ou
WHERE u.user_uuid = sqlc.arg(current_user_uuid)::UUID
  AND ou.user_uuid = sqlc.arg(other_user_uuid)::UUID
RETURNING *;

-- name: GetDirectMessagesConversations :many
WITH cur_user AS (
    SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID
)
SELECT
    c.conversation_id,
    c.last_message_at,
    other_user.user_uuid AS other_user_uuid,
    other_user.user_email as other_user_email,
    dm.dm_content AS last_message_content,
    dm.dm_message_type AS last_message_type,
    (
        SELECT count(*)
        FROM direct_messages dm2
        WHERE dm2.receiver_id = (SELECT user_id FROM cur_user)
            AND dm2.sender_id = other_user.user_id
            AND (cp.last_read_message_id IS NULL OR dm2.dm_id > cp.last_read_message_id)
            AND dm2.dm_deleted_at IS NULL
    ) AS unread_count
FROM conversations c
JOIN conversation_participants cp
    ON cp.conversation_id = c.conversation_id
    AND cp.user_id = (SELECT user_id FROM cur_user)
    AND cp.left_at IS NULL
JOIN conversation_participants cp_other
    ON cp_other.conversation_id = c.conversation_id
    AND cp_other.user_id != (SELECT user_id FROM current_user)
    AND cp_other.left_at IS NULL
JOIN users other_user
    ON other_user.user_id = cp_other.user_id
LEFT JOIN direct_messages dm
    ON dm.dm_id = c.last_message_id
    AND dm.dm_deleted_at IS NULL
WHERE c.conversation_type = 1
ORDER BY c.last_message_at DESC NULLS LAST
LIMIT sqlc.arg(limitArg) OFFSET sqlc.arg(offsetArg);

-- name: GetDirectMessageWithUser :many
WITH user_ids AS(
    SELECT
        u.user_id AS current_user_id,
        ou.user_id AS other_user_id
    FROM users u, users ou
    WHERE
        u.user_uuid = sqlc.arg(current_user_uuid)::UUID
        AND ou.user_uuid = sqlc.arg(other_user_uuid)::UUID
)
SELECT
    dm.dm_id,
    CASE
        WHEN dm.sender_id = uid.current_user_id THEN TRUE ELSE FALSE
    END AS is_mine,
    dm.dm_content,
    dm.dm_attachments,
    dm.dm_message_type,
    dm.dm_read_at,
    dm.dm_created_at
FROM direct_messages dm
JOIN user_ids uid
    ON (dm.sender_id = uid.current_user_id AND dm.receiver_id = uid.other_user_id)
    OR (dm.sender_id = uid.other_user_id AND dm.receiver_id = uid.current_user_id)
WHERE dm.dm_deleted_at IS NULL
ORDER BY dm.dm_created_at DESC
LIMIT sqlc.arg(limitArg) OFFSET sqlc.arg(offsetArg);

-- name: EditDirectMessage :one
UPDATE direct_messages
SET
    dm_content = sqlc.arg(dm_content)
WHERE dm_id = sqlc.arg(dm_id)
    AND sender_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(current_user_uuid)::UUID)
    AND dm.dm_deleted_at IS NULL
    AND dm.dm_read_at IS NULL
    AND dm_message_type = 1
RETURNING *;

-- name: SoftDeleteDirectMessage :one
UPDATE direct_messages
SET
    dm_deleted_at = now()
WHERE dm_id = sqlc.arg(dm_id)
    AND sender_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(current_user_uuid)::UUID)
    AND dm.dm_deleted_at IS NULL
RETURNING *;

-- name: ReadDirectMessage :one
UPDATE direct_messages
SET
    dm_read_at = now()
WHERE dm_id = sqlc.arg(dm_id)
    AND receiver_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(current_user_uuid)::UUID)
    AND dm.dm_deleted_at IS NULL
    AND dm.dm_read_at IS NULL
RETURNING *;

-- name: ReadAllDirectMessages :one
WITH updated AS (
    UPDATE direct_messages
    SET dm_read_at = now()
    WHERE sender_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(receiver_uuid)::UUID)
      AND receiver_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(current_user_uuid)::UUID)
      AND dm_deleted_at IS NULL
      AND dm_read_at IS NULL
    RETURNING dm_id
)
SELECT count(*) AS read_count FROM updated;
-- name: GetUnreadDM :one
WITH unread_msgs AS (
    SELECT dm_id, sender_id
    FROM direct_messages
    WHERE receiver_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(current_user_uuid)::UUID)
      AND dm_deleted_at IS NULL
      AND dm_read_at IS NULL
),
total AS (
    SELECT COUNT(*) AS total_unread FROM unread_msgs
),
by_user AS (
    SELECT COALESCE(
        jsonb_agg(jsonb_build_object('user_id', sender_id, 'count', c) ORDER BY user_id),
        '[]'::jsonb
    ) AS unread_by_user
    FROM (
        SELECT sender_id AS user_id, COUNT(*) AS c
        FROM unread_msgs
        GROUP BY sender_id
    ) t
)
SELECT total.total_unread, by_user.unread_by_user
FROM total, by_user;

-- name: SearchDirectMessagesWithUser :many
WITH user_ids AS (
    SELECT
        cu.user_id AS current_user_id,
        ou.user_id AS other_user_id
    FROM users cu, users ou
    WHERE cu.user_uuid = sqlc.arg(current_user_uuid)::UUID
      AND ou.user_uuid = sqlc.arg(other_user_uuid)::UUID
)
SELECT
    dm.dm_id,
    dm.sender_id,
    dm.receiver_id,
    dm.dm_content,
    dm.dm_attachments,
    dm.dm_message_type,
    dm.dm_read_at,
    dm.dm_created_at
FROM direct_messages dm
JOIN user_ids uid
    ON (dm.sender_id = uid.current_user_id AND dm.receiver_id = uid.other_user_id)
    OR (dm.sender_id = uid.other_user_id AND dm.receiver_id = uid.current_user_id)
WHERE dm.dm_content ILIKE '%' || sqlc.arg(search) || '%'
  AND dm.dm_deleted_at IS NULL;
