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
