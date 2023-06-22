-- name: GetSpace :one
SELECT *
FROM spaces
WHERE id = $1
LIMIT 1;

-- name: GetPost :one
SELECT *
FROM posts
WHERE id = $1
LIMIT 1;

-- name: ListPosts :many
SELECT *
FROM posts
WHERE space_id = $1;

-- name: CreateSpace :exec
INSERT INTO spaces (name, description, parent_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteSpace :exec
DELETE
FROM spaces
WHERE id = $1;