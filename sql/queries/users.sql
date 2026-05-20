-- name: CreateUser :one
INSERT INTO users (id, name, api_key, created_at, updated_at)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT id, name, api_key, created_at, updated_at
FROM users
WHERE api_key = $1;