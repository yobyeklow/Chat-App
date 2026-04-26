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
    SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)
)
    AND gm.left_at IS NULL
    AND gm.group_deleted_at IS NULL
ORDER BY gm.jointed_at DESC
LIMIT sqlc.arg(limitArg) OFFSET sqlc.arg(offsetArg);

-- name: GetGroupByUUID :one
SELECT *
FROM groups g
WHERE g.group_uuid = sqlc.arg(group_uuid)::UUID;
-- name: SoftDeleteGroup :one
UPDATE groups
SET
    group_deleted_at = now()
WHERE group_uuid = sqlc.arg(group_uuid)::UUID and group_deleted_at IS NULL
RETURNING *;
-- name: HardDeleteGroup :one
DELETE FROM groups
WHERE group_uuid = sqlc.arg(group_uuid)::UUID and group_deleted_at IS NULL
RETURNING *;
