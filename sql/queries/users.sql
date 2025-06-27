-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email) 
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: DeleteUsers :exec
DELETE FROM users *;

-- name: GetUser :one
SELECT id FROM users
WHERE id = $1;