-- name: GetGroupMembers :many
SELECT
    g.group_uuid,
    us.user_uuid,
    gm.member_role,
    gm.jointed_at
FROM group_members gm
INNER JOIN users us ON us.user_id = gm.user_id
INNER JOIN groups g ON g.group_id = gm.group_id
WHERE
    g.group_uuid = sqlc.arg(group_uuid)::UUID
    AND gm. left_at IS NULL
ORDER BY gm.jointed_at DESC
LIMIT sqlc.arg(limitArg) OFFSET sqlc.arg(offsetArg);

-- name: AddMemberToGroup :one
INSERT INTO group_members(group_id, user_id, member_role, jointed_at)
SELECT
    (SELECT group_id FROM groups WHERE group_uuid = sqlc.arg(group_uuid)::UUID),
    (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID),
    sqlc.arg(member_role),
    now()
RETURNING *;

-- name: UpdateMemberRole :one
UPDATE group_members
SET
    member_role = sqlc.arg(member_role)
WHERE
    group_id = (SELECT group_id FROM groups WHERE group_uuid = sqlc.arg(group_uuid)::UUID)
    AND user_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID)
    AND left_at IS NULL
RETURNING *;
-- name: RemoveMember :one
UPDATE group_members
SET
    left_at = now()
WHERE
    group_id = (SELECT group_id FROM groups WHERE group_uuid = sqlc.arg(group_uuid)::UUID)
    AND user_id = (SELECT user_id FROM users WHERE user_uuid = sqlc.arg(user_uuid)::UUID)
    AND left_at IS NULL
RETURNING *;
-- name: GetMemberInfo :one
SELECT
    us.user_uuid,
    us.user_email,
    g.group_uuid,
    gm.member_role,
    gm.jointed_at
FROM group_members gm
INNER JOIN users us ON us.user_id = gm.user_id
INNER JOIN groups g ON g.group_id = gm.group_id
WHERE
    us.user_uuid = sqlc.arg(user_uuid)::UUID
    AND g.group_uuid = sqlc.arg(group_uuid)::UUID
    AND gm.left_at IS NULL;
