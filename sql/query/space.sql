-- name: GetSpace :one
SELECT sqlc.embed(s),
       small_picture.url         as space_small_pic_url,
       small_picture.width       as space_small_pic_width,
       small_picture.height      as space_small_pic_height,
       background_picture.url    as background_picture_url,
       background_picture.width  as background_picture_width,
       background_picture.height as background_picture_height
FROM spaces s
         left join pictures small_picture on s.small_picture_id = small_picture.id
         left join pictures background_picture on s.background_picture_id = background_picture.id
WHERE s.id = $1
LIMIT 1;

-- name: GetSpacesOfParent :many
SELECT sqlc.embed(s),
       small_picture.url         as space_small_pic_url,
       small_picture.width       as space_small_pic_width,
       small_picture.height      as space_small_pic_height,
       background_picture.url    as background_picture_url,
       background_picture.width  as background_picture_width,
       background_picture.height as background_picture_height
FROM spaces s
         left join pictures small_picture on s.small_picture_id = small_picture.id
         left join pictures background_picture on s.background_picture_id = background_picture.id
WHERE s.parent_id = $1;

-- name: GetSpaceByName :one
SELECT sqlc.embed(s),
       small_picture.url         as space_small_pic_url,
       small_picture.width       as space_small_pic_width,
       small_picture.height      as space_small_pic_height,
       background_picture.url    as background_picture_url,
       background_picture.width  as background_picture_width,
       background_picture.height as background_picture_height
FROM spaces s
         left join pictures small_picture on s.small_picture_id = small_picture.id
         left join pictures background_picture on s.background_picture_id = background_picture.id
WHERE s.name = $1
  AND s.parent_id = $2
LIMIT 1;

-- name: GetUserSpaces :many
SELECT sqlc.embed(sp),
       small_picture.url         as space_small_pic_url,
       small_picture.width       as space_small_pic_width,
       small_picture.height      as space_small_pic_height,
       background_picture.url    as background_picture_url,
       background_picture.width  as background_picture_width,
       background_picture.height as background_picture_height
FROM subscriptions su
         JOIN spaces sp ON su.space_id = sp.id
         left join pictures small_picture on s.small_picture_id = small_picture.id
         left join pictures background_picture on s.background_picture_id = background_picture.id
WHERE su.user_id = $1
  AND su.is_deleted = false;

-- name: SearchSpace :many
SELECT sqlc.embed(s),
       small_picture.url         as space_small_pic_url,
       small_picture.width       as space_small_pic_width,
       small_picture.height      as space_small_pic_height,
       background_picture.url    as background_picture_url,
       background_picture.width  as background_picture_width,
       background_picture.height as background_picture_height
FROM spaces s
         left join pictures small_picture on s.small_picture_id = small_picture.id
         left join pictures background_picture on s.background_picture_id = background_picture.id
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

