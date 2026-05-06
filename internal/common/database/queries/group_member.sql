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

-- name: AddMemberToGroup :exec
WITH target_group AS(
    SELECT group_id
    FROM groups
    WHERE group_uuid = sqlc.arg(group_uuid)::UUID
        AND group_deleted_at IS NULL
),
target_user AS(
    SELECT user_id
    FROM users
    WHERE user_uuid = sqlc.arg(user_uuid)::UUID
        AND user_deleted_at IS NULL
        AND user_status = 1
),
insert_member AS(
    INSERT INTO group_members (group_id, user_id, member_role)
    SELECT tg.group_id, tu.user_id, 1
    FROM target_group tg, target_user tu
    ON CONFLICT (group_id,user_id) WHERE left_at IS NULL DO UPDATE
        SET left_at = NULL,
            jointed_at = now(),
            member_role = EXCLUDED.member_role
    RETURNING group_member_id, user_id
),
target_conversation AS(
    SELECT conversation_id FROM conversations
    WHERE reference_id = (SELECT group_id FROM target_group)
        AND conversation_type = 2
)
INSERT INTO conversation_participants (conversation_id, user_id)
SELECT tc.conversation_id, im.user_id
FROM target_conversation tc, insert_member im
ON CONFLICT (conversation_id,user_id) WHERE left_at IS NULL DO UPDATE
    SET left_at = NULL,
        joined_at = now();

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
