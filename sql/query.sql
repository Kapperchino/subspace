-- name: GetSpace :one
SELECT *
FROM spaces
WHERE id = $1
LIMIT 1;

-- name: GetSpacesOfParent :many
SELECT *
FROM spaces
WHERE parent_id = $1;

-- name: GetSpaceByName :one
SELECT *
FROM spaces
WHERE name = $1
LIMIT 1;

-- name: CreateSpace :exec
INSERT INTO spaces (name, description, parent_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteSpace :exec
DELETE
FROM spaces
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (password, email, display_name, bio)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserFromEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: CreatePost :one
INSERT INTO posts (space_id, poster_id, topic, content)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPost :one
SELECT *
FROM posts
WHERE id = $1
LIMIT 1;

-- name: GetPostsForSpace :many
SELECT *
FROM posts
WHERE space_id = $1;