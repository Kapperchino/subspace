-- name: CreatePost :one
INSERT INTO posts (space_id, poster_id, topic, body, content_type, link)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: SearchPostPopular :many
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v on v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id != 1
  AND current_timestamp - p.created
    < make_interval(days => $2)
ORDER BY ts_rank(p.ts, plainto_tsquery('english', $3)) DESC, up_votes DESC
LIMIT 10 OFFSET $4;

-- name: SearchPostLatest :many
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v on v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id != 1
  AND current_timestamp - p.created
    < make_interval(days => $2)
ORDER BY ts_rank(p.ts, plainto_tsquery('english', $3)) DESC, p.created DESC
LIMIT 10 OFFSET $4;

-- name: GetPost :one
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v on v.user_id = $2 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id = $1
LIMIT 1;

-- name: GetPostsForSpaceLatestByName :many
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v on v.user_id = $3 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.space_name = $1
  AND p.parent_id = $2
  AND p.id != 1
  AND current_timestamp - p.created < make_interval(days => $4)
ORDER BY p.created DESC
LIMIT 10 OFFSET $5;

-- name: GetPostsForSpaceLatest :many
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v on v.user_id = $2 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.space_id = $1
  AND p.id != 1
  AND current_timestamp - p.created < make_interval(days => $3)
ORDER BY p.created DESC
LIMIT 10 OFFSET $4;

-- name: GetPostsForHomeLatest :many
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v on v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id != 1
  AND current_timestamp - p.created < make_interval(days => $2)
ORDER BY p.created DESC
LIMIT 10 OFFSET $3;

-- name: GetPostsForSpacePopularByName :many
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v on v.user_id = $3 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.space_name = $1
  AND p.parent_id = $2
  AND p.id != 1
  AND current_timestamp - p.created < make_interval(days => $4)
ORDER BY up_votes DESC, p.id DESC
LIMIT 10 OFFSET $5;

-- name: GetPostsForSpacePopular :many
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v
                   on v.user_id = $2 and p.id = v.post_or_comment_id and
                      v.vote_type = 'post'
WHERE p.space_id = $1
  AND p.id != 1
  AND current_timestamp - p.created
    < make_interval(days => $3)
ORDER BY up_votes DESC, p.id DESC
LIMIT 10 OFFSET $4;

-- name: GetPostsForHomePopular :many
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v
                   on v.user_id = $1 and p.id = v.post_or_comment_id and
                      v.vote_type = 'post'
WHERE p.id != 1
  AND current_timestamp - p.created
    < make_interval(days => $2)
ORDER BY up_votes DESC, p.id DESC
LIMIT 10 OFFSET $3;

-- name: GetPostsForUserLatest :many
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v
                   on v.user_id = $1 and p.id = v.post_or_comment_id and
                      v.vote_type = 'post'
WHERE p.poster_id = $1
  AND p.id != 1
  AND current_timestamp - p.created < make_interval(days => $2)
ORDER BY p.created DESC
LIMIT 10 OFFSET $3;

-- name: GetPostsForUserPopular :many
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v
                   on v.user_id = $1 and p.id = v.post_or_comment_id and
                      v.vote_type = 'post'
WHERE p.poster_id = $1
  AND p.id != 1
  AND current_timestamp - p.created < make_interval(days => $2)
ORDER BY up_votes DESC, p.id DESC
LIMIT 10 OFFSET $3;

-- name: GetPostsForUserSubscriptionLatest :many
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v
                   on v.user_id = $1 and p.id = v.post_or_comment_id and
                      v.vote_type = 'post'
         join subscriptions su on su.space_id = p.space_id and su.user_id = $1 and su.is_deleted = false
WHERE p.id != 1
  AND p.space_id != 1
  AND current_timestamp - p.created < make_interval(days => $2)
ORDER BY p.created DESC
LIMIT 10 OFFSET $3;

-- name: GetPostsForUserSubscriptionPopular :many
SELECT sqlc.embed(p), v.*
from posts_view p
         left join votes v
                   on v.user_id = $1 and p.id = v.post_or_comment_id and
                      v.vote_type = 'post'
         join subscriptions su on su.space_id = p.space_id and su.user_id = $1 and su.is_deleted = false
WHERE p.id != 1
  AND p.space_id != 1
  AND current_timestamp - p.created < make_interval(days => $2)
ORDER BY up_votes DESC, p.id DESC
LIMIT 10 OFFSET $3;

-- name: GetPoster :one
SELECT u.*
from posts p
         join users u on p.poster_id = u.id
where p.id = $1;