-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;