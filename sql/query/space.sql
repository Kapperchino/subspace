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
  AND parent_id = $2
LIMIT 1;

-- name: GetUserSpaces :many
SELECT sp.*
FROM subscriptions su
         JOIN spaces sp ON su.space_id = sp.id
WHERE su.user_id = $1
  AND su.is_deleted = false;

-- name: SearchSpace :many
SELECT sqlc.embed(s), sqlc.embed(p)
FROM spaces s
         LEFT JOIN pictures p on s.small_picture_id = p.id
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

