-- name: CreateUser :one
INSERT INTO users (password, email, display_name, bio)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT sqlc.embed(u), p.*
FROM users u
         left join pictures p on u.picture_id = p.id
WHERE u.id = $1
LIMIT 1;

-- name: GetUserFromEmail :one
SELECT sqlc.embed(u), p.*
FROM users u
         left join pictures p on u.picture_id = p.id
WHERE u.email = $1
LIMIT 1;