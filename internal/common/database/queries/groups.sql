-- name: CreateGroup :one
INSERT INTO groups(
    group_name
) VALUES (sqlc.arg(group_name))
ON CONFLICT (group_name) DO NOTHING
RETURNING *;

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
WHERE gm.member_id = (
    SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID
)
    AND gm.left_at IS NULL
    AND g.group_deleted_at IS NULL
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
WHERE gm.member_id = (
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
UPDATE group_members
SET left_at = now()
WHERE member_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID)
    AND group_id = (SELECT group_id FROM groups WHERE group_uuid = sqlc.arg(group_uuid)::UUID)
    AND left_at IS NULL;

-- name: UpdateGroup :one
UPDATE groups
SET group_name = sqlc.arg(group_name)
WHERE group_uuid = sqlc.arg(group_uuid)::UUID
    AND group_deleted_at IS NULL
RETURNING *;

-- name: HardDeleteGroup :exec
DELETE FROM groups
WHERE group_uuid = sqlc.arg(group_uuid)::UUID and group_deleted_at IS NULL;

-- name: GetMemberRole :one
SELECT member_role
FROM group_members
WHERE member_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID)
    AND group_id = (SELECT group_id FROM groups WHERE group_uuid = sqlc.arg(group_uuid)::UUID)
    AND left_at IS NULL;
