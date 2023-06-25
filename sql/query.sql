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
LIMIT 1;

-- name: CreateSpace :one
INSERT INTO spaces (name, description, parent_id, picture)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteSpace :exec
DELETE
FROM spaces
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (password, email, display_name, bio)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserFromEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: CreatePost :one
INSERT INTO posts (space_id, poster_id, topic, body, content, content_type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPost :one
SELECT p.*, u.display_name, s.picture as space_picture
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
WHERE p.id = $1
LIMIT 1;

-- name: GetPostsForSpace :many
SELECT p.*, u.display_name, s.picture as space_picture
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
WHERE space_id = $1;

-- name: CreateVote :one
INSERT INTO votes (user_id, post_or_comment_id, vote_type, is_up_vote)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: SetVote :exec
UPDATE votes
set is_up_vote = $1
where id = $2;

-- name: DeleteVote :exec
UPDATE votes
set is_deleted = true
where id = $1;

-- name: GetVoteForPostOrCommentForUser :one
SELECT *
FROM votes
WHERE user_id = $1
  AND post_or_comment_id = $2
  AND is_deleted = false
LIMIT 1;

-- name: GetVotesForPostOrComment :many
SELECT *
FROM votes
WHERE post_or_comment_id = $1
  AND is_deleted = false
LIMIT 1;

-- name: UpVotePost :exec
UPDATE posts
SET up_votes = up_votes + 1
WHERE poster_id = $1;

-- name: DeUpVotePost :exec
UPDATE posts
SET up_votes = up_votes - 1
WHERE poster_id = $1;

-- name: DownVotePost :exec
UPDATE posts
SET down_votes = down_votes + 1
WHERE poster_id = $1;

-- name: DeDownVotePost :exec
UPDATE posts
SET down_votes = down_votes - 1
WHERE poster_id = $1;

-- name: UpVoteComment :exec
UPDATE comments
SET up_votes = up_votes + 1
WHERE poster_id = $1;

-- name: DeUpVoteComment :exec
UPDATE comments
SET up_votes = up_votes - 1
WHERE poster_id = $1;

-- name: DownVoteComment :exec
UPDATE comments
SET down_votes = down_votes + 1
WHERE poster_id = $1;

-- name: DeDownVoteComment :exec
UPDATE comments
SET down_votes = down_votes - 1
WHERE poster_id = $1;
