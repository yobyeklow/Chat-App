-- name: CreateGroupMember :one
INSERT INTO group_members(
    group_id,
    user_id,
    member_role
) VALUES ($1,$2,$3)
RETURNING *;
