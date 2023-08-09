-- name: CreatePicture :one
INSERT INTO pictures (url, width, height)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreatePictureRelation :one
INSERT INTO picture_releations (picture_id, post_id, comment_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPicturesForPost :many
SELECT p.*
FROM pictures p
         join picture_releations pr on p.id = pr.picture_id
where pr.post_id = $1;

-- name: GetPicture :one
SELECT p.*
FROM pictures p
where p.id = $1;

-- name: GetPicturesForComment :many
SELECT p.*
FROM pictures p
         join picture_releations pr on p.id = pr.picture_id
where pr.post_id = $1;

-- name: GetPicturesForPosts :many
SELECT p.*
FROM pictures p
         join picture_releations pr on p.id = pr.picture_id
where pr.post_id IN ($1::bigint[]);