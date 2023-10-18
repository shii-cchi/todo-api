-- name: CreateUser :one
INSERT INTO users (name)
VALUES ($1)
RETURNING *;

-- name: CheckApiKey :one
SELECT id FROM users
WHERE api_key = $1;