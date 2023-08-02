-- name: CreatePicture :one
INSERT INTO pictures (post_id, comment_id, url, width, height)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;