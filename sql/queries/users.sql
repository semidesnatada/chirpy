-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET updated_at = NOW(), email = $2, hashed_password = $3
WHERE id = $1
RETURNING *;

-- name: GetHashedPassword :one
SELECT hashed_password
FROM users
WHERE email = $1;

-- name: GetUser :one
SELECT *
FROM users
WHERE email = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: UpgradeUser :one
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING *;