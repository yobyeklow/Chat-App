-- name: CreateUser :one
INSERT INTO users(
    user_email,
    user_password,
    user_status,
    user_role
)VALUES(
    $1,$2,$3,$4
) RETURNING *;
-- name: FindUserByEmail :one
SELECT *
FROM users
WHERE user_email = sqlc.arg(user_email)::TEXT and user_deleted_at IS NULL;
-- name: FindUserByUUID :one
SELECT *
FROM users
WHERE user_uuid = sqlc.arg(user_uuid)::UUID and user_deleted_at IS NULL;
-- name: SoftDelete :one
UPDATE users
SET
    user_deleted_at = now(),
    user_status = 2
WHERE
    user_uuid = sqlc.arg(user_uuid)::UUID AND user_deleted_at IS NULL
RETURNING *;
-- name: HardDelete :one
DELETE FROM users
WHERE
    user_uuid = sqlc.arg(user_uuid)::UUID AND user_deleted_at IS NOT NULL
RETURNING *;
-- name: RestoreUser :one
UPDATE users
SET
    user_deleted_at = NULL
WHERE
    user_uuid = sqlc.arg(user_uuid)::uuid AND user_deleted_at IS NOT NULL
RETURNING *;
