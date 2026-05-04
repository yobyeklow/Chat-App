-- name: CreateGroup :one
INSERT INTO groups(
    group_name
) VALUES (sqlc.arg(group_name))
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
UPDATE group_members
SET left_at = now()
WHERE user_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID)
    AND group_id = (SELECT group_id FROM groups WHERE group_uuid = sqlc.arg(group_uuid)::UUID)
    AND left_at IS NULL;

-- name: UpdateGroup :one
UPDATE groups
SET group_name = sqlc.arg(group_name)
WHERE group_uuid = sqlc.arg(group_uuid)::UUID
    AND group_deleted_at IS NULL
RETURNING *;

-- name: HardDeleteGroup :one
DELETE FROM groups
WHERE group_uuid = sqlc.arg(group_uuid)::UUID and group_deleted_at IS NOT NULL RETURNING *;

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
