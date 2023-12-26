-- name: GetSpace :one
SELECT *
FROM spaces_view s
WHERE s.id = $1
LIMIT 1;

-- name: GetSpacesOfParent :many
SELECT *
from spaces_view s
WHERE s.parent_id = $1;

-- name: GetSpaceByName :one
SELECT *
FROM spaces_view s
WHERE s.name = $1
  AND s.parent_id = $2
LIMIT 1;

-- name: GetUserSpaces :many
SELECT sp.*
FROM subscriptions su
         JOIN spaces_view sp ON su.space_id = sp.id
WHERE su.user_id = $1
  AND su.is_deleted = false;

-- name: GetPopularSpaces :many
SELECT *
FROM spaces_view s
WHERE s.id != 1
ORDER BY s.sub_count DESC;

-- name: GetNewSpaces :many
SELECT *
FROM spaces_view s
WHERE s.id != 1
ORDER BY s.created DESC;

-- name: SearchSpace :many
SELECT *
FROM spaces_view s
WHERE s.id != 1
ORDER BY ts_rank(ts, plainto_tsquery('english', $1)) DESC;

-- name: CreateSpace :one
INSERT INTO spaces (name, description, parent_id, small_picture_id, background_picture_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteSpace :exec
DELETE
FROM spaces
WHERE id = $1;

-- name: SpacePrefixSearch :many
select sqlc.embed(s), similarity(name, $1) as sml
from spaces_view s
where s.id != 1
order by sml desc
limit 10;

