-- name: CreateUser :one
INSERT INTO users (id, name, password, api_key, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE name = $1;

-- name: GetUserByName :one
SELECT * FROM users
WHERE name = $1;