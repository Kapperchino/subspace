-- name: CreatePicture :one
INSERT INTO pictures (url, width, height)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPicturesForPost :many
SELECT p.*
FROM pictures p
         join pictureRelations pr on p.id = pr.picture_id
where pr.post_id = $1;

-- name: GetPicturesForPosts :many
SELECT p.*
FROM pictures p
         join pictureRelations pr on p.id = pr.picture_id
where pr.post_id IN ($1::bigint[]);