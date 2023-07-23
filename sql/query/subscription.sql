-- name: CreateSubscription :one
INSERT INTO subscriptions (user_id, space_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetSubscription :one
SELECT *
FROM subscriptions
where user_id = $1
  AND space_id = $2;

-- name: GetActiveSubscription :one
SELECT *
FROM subscriptions
where user_id = $1
  AND space_id = $2
  AND is_deleted = false;


-- name: GetSubscriptionsForUser :many
SELECT *
FROM subscriptions
where user_id = $1;

-- name: RefreshSubscription :one
UPDATE subscriptions
SET is_deleted = false
where id = $1
RETURNING *;

-- name: DeleteSubscription :exec
UPDATE subscriptions
set is_deleted = true
where space_id = $1
  AND user_id = $2;


