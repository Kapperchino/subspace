-- name: CreateUserAddress :one
INSERT INTO user_addresses (from_user_id, to_user_id, post_id)
VALUES ($1, $2, $3)
RETURNING *;