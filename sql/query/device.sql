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