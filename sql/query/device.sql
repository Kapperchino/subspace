-- name: CreateDevice :one
INSERT INTO devices (user_id, registration, device_info)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDeviceByDeviceInfo :one
SELECT *
FROM devices
WHERE device_info = $1
  AND user_id = $2;

-- name: UpdateDevice :one
UPDATE devices
set registration = $1
where device_info = $2
  AND user_id = $3
RETURNING *;

-- name: GetDevicesFromCommentId :many
SELECT *
FROM devices
where user_id = (SELECT c.poster_id
                 FROM comments c
                 where c.id = $1);

-- name: GetDevicesFromPostId :many
SELECT *
FROM devices
where user_id = (SELECT p.poster_id
                 FROM posts p
                 where p.id = $1);