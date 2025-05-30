-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    token,
    created_at,
    updated_at,
    user_id,
    expires_at,
    revoked_at
    )
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4
)
RETURNING *;

-- name: CheckRefreshTokenExistsAndIsValid :one
SELECT EXISTS (SELECT 1
FROM refresh_tokens
WHERE $1 = token AND revoked_at IS NULL
LIMIT 1);

-- name: GetUserFromAccessToken :one
SELECT user_id
FROM refresh_tokens
WHERE token = $1
LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;