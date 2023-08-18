-- name: CreateUser :one
INSERT INTO users (password, email, display_name, bio, address)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUserBio :exec
UPDATE users
set bio = $1
where id = $2
RETURNING *;

-- name: UpdatePicture :exec
UPDATE users
SET picture_id = $1
WHERE id = $2;

-- name: GetUser :one
SELECT sqlc.embed(u), p.*
FROM users u
         left join pictures p on u.picture_id = p.id
WHERE u.id = $1
LIMIT 1;

-- name: GetUserFromAddress :one
SELECT sqlc.embed(u), p.*
FROM users u
         left join pictures p on u.picture_id = p.id
WHERE u.address = $1
LIMIT 1;

-- name: GetUserFromEmail :one
SELECT sqlc.embed(u), p.*
FROM users u
         left join pictures p on u.picture_id = p.id
WHERE u.email = $1
LIMIT 1;

-- name: SearchUsers :many
SELECT sqlc.embed(u), p.*
from users u
         left join pictures p on u.picture_id = p.id
where u.is_deleted = false
  AND u.id != 1
ORDER BY ts_rank(u.ts, plainto_tsquery('english', $1)) DESC
LIMIT 100;