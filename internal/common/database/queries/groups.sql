-- name: CreateGroup :one
WITH new_group AS (
    INSERT INTO groups (group_name)
    VALUES (sqlc.arg(group_name))
    RETURNING *
),
new_conversation AS (
    INSERT INTO conversations (conversation_type, reference_id)
    SELECT 2, ng.group_id
    FROM new_group ng
    RETURNING *
),
creator_user AS(
    SELECT user_id FROM users
    WHERE user_uuid =  sqlc.arg(user_uuid)::UUID
        AND user_deleted_at IS NULL
        AND user_status = 1
),
new_member AS(
    INSERT INTO group_members(group_id,user_id,member_role)
    SELECT  ng.group_id,cu.user_id,3
    FROM new_group ng,creator_user cu
    ON CONFLICT(group_id,user_id) WHERE left_at IS NULL DO NOTHING
    RETURNING *
),
new_participant AS(
    INSERT INTO conversation_participants (conversation_id, user_id)
    SELECT nc.conversation_id, nm.user_id
    FROM new_conversation nc, new_member nm
    ON CONFLICT (conversation_id,user_id) WHERE left_at IS NULL DO NOTHING
    RETURNING *
)
SELECT
    ng.group_id,
    ng.group_uuid,
    ng.group_name,
    ng.group_created_at,
    ng.group_updated_at,
    nc.conversation_id,
    nc.conversation_uuid
FROM new_group ng
JOIN new_conversation nc ON true;
-- name: GetAllGroups :many
SELECT DISTINCT
    g.group_uuid,
    g.group_name,
    g.group_created_at,
    g.group_updated_at,
    gm.member_role,
    gm.jointed_at
FROM groups g
INNER JOIN group_members gm on g.group_id = gm.group_id
WHERE gm.user_id = (
    SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID
)
    AND gm.left_at IS NULL
    AND g.group_deleted_at IS NULL
    AND (
            sqlc.arg(search) = ''
            OR to_tsvector('english', COALESCE(g.group_name, '')) @@
            plainto_tsquery('english', sqlc.arg(search)::TEXT)
        )
ORDER BY gm.jointed_at DESC
LIMIT sqlc.arg(limitArg) OFFSET sqlc.arg(offsetArg);

-- name: GetGroupByUUID :one
SELECT
    g.group_uuid,
    g.group_name,
    g.group_created_at,
    g.group_updated_at,
    gm.member_role,
    gm.jointed_at
FROM groups g
INNER JOIN  group_members gm ON g.group_id = gm.group_id
WHERE gm.user_id = (
    SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID
)
    AND gm.left_at IS NULL
    AND g.group_deleted_at IS NULL
    AND g.group_uuid = sqlc.arg(group_uuid)::UUID;

-- name: SoftDeleteGroup :one
UPDATE groups
SET group_deleted_at = now()
WHERE group_uuid = sqlc.arg(group_uuid)::UUID
    AND group_deleted_at IS NULL
RETURNING *;
-- name: LeaveGroup :exec
WITH uid AS (
        SELECT user_id FROM users
        WHERE user_uuid = sqlc.arg(user_uuid)::UUID
),
gid AS (
    SELECT group_id FROM groups
    WHERE group_uuid = sqlc.arg(group_uuid)::UUID
        AND group_deleted_at IS NULL
),
cid AS (
    SELECT conversation_id FROM conversations
    WHERE reference_id = (SELECT group_id FROM gid)
        AND conversation_type = 2   -- 2 = Group chat
),
leave_member AS (
    UPDATE group_members
    SET left_at = now()
    WHERE user_id = (SELECT user_id FROM uid)
        AND group_id = (SELECT group_id FROM gid)
        AND left_at IS NULL
)
UPDATE conversation_participants
SET left_at = now()
WHERE conversation_id = (SELECT conversation_id FROM cid)
  AND user_id = (SELECT user_id FROM uid)
  AND left_at IS NULL;

-- name: UpdateGroup :one
UPDATE groups
SET group_name = sqlc.arg(group_name)
WHERE group_uuid = sqlc.arg(group_uuid)::UUID
    AND group_deleted_at IS NULL
RETURNING *;

-- name: HardDeleteGroup :exec
DELETE FROM groups
WHERE group_uuid = sqlc.arg(group_uuid)::UUID and group_deleted_at IS NOT NULL;

-- name: GetMemberRole :one
SELECT member_role
FROM group_members
WHERE user_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID)
    AND group_id = (SELECT group_id FROM groups WHERE group_uuid = sqlc.arg(group_uuid)::UUID)
    AND left_at IS NULL;
-- name: CountRecords :one
SELECT count(*)
FROM groups
WHERE (
    sqlc.narg(deleted)::bool IS NULL
    OR (group_deleted_at IS NOT NULL AND sqlc.narg(deleted)::bool IS TRUE)
    OR (group_deleted_at IS NULL AND sqlc.narg(deleted)::bool IS FALSE)
) AND (
    sqlc.narg(search)::TEXT IS NULL
    OR sqlc.narg(search)::TEXT = ''
    OR group_name ILIKE '%' || sqlc.narg(search) || '%');
