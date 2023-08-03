-- name: GetSpace :one
SELECT sqlc.embed(s), sqlc.embed(small_picture), sqlc.embed(background_picture)
FROM spaces s
         left join pictures small_picture on s.small_picture_id
         left join pictures background_picture on s.background_picture_id
WHERE s.id = $1
LIMIT 1;

-- name: GetSpacesOfParent :many
SELECT sqlc.embed(s), sqlc.embed(small_picture), sqlc.embed(background_picture)
FROM spaces s
         left join pictures small_picture on s.small_picture_id
         left join pictures background_picture on s.background_picture_id
WHERE s.parent_id = $1;

-- name: GetSpaceByName :one
SELECT sqlc.embed(s), sqlc.embed(small_picture), sqlc.embed(background_picture)
FROM spaces s
         left join pictures small_picture on s.small_picture_id
         left join pictures background_picture on s.background_picture_id
WHERE name = $1
  AND parent_id = $2
LIMIT 1;

-- name: GetUserSpaces :many
SELECT sqlc.embed(sp), sqlc.embed(small_picture), sqlc.embed(background_picture)
FROM subscriptions su
         JOIN spaces sp ON su.space_id = sp.id
         left join pictures small_picture on s.small_picture_id
         left join pictures background_picture on s.background_picture_id
WHERE su.user_id = $1
  AND su.is_deleted = false;

-- name: SearchSpace :many
SELECT sqlc.embed(s), sqlc.embed(small_picture), sqlc.embed(background_picture)
FROM spaces s
         left join pictures small_picture on s.small_picture_id
         left join pictures background_picture on s.background_picture_id
WHERE id != 1
ORDER BY ts_rank(ts, plainto_tsquery('english', $1)) DESC;

-- name: CreateSpace :one
INSERT INTO spaces (name, description, parent_id, small_picture_id, background_picture_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteSpace :exec
DELETE
FROM spaces
WHERE id = $1;

