-- name: CreatePost :one
INSERT INTO posts (space_id, poster_id, topic, body, content, content_type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: SearchPost :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       s.parent_id,
       s.name                       as space_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'post') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'post') AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id != 1
  AND current_timestamp - p.created
    < make_interval(days => $2)
ORDER BY ts_rank(p.ts, plainto_tsquery('english', $3)) DESC, up_votes DESC;

-- name: GetPost :one
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       s.parent_id,
       s.name                       as space_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'post') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'post') AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on v.user_id = $2 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id = $1
LIMIT 1;

-- name: GetPostsForSpaceLatestByName :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       s.parent_id,
       s.name                       as space_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'post') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'post') AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on v.user_id = $3 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE s.name = $1
  AND s.parent_id = $2
  AND p.id != 1
  AND current_timestamp - p.created < make_interval(days => $4)
ORDER BY p.created DESC;

-- name: GetPostsForSpaceLatest :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       s.parent_id,
       s.name                       as space_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'post') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'post') AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on v.user_id = $2 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE space_id = $1
  AND p.id != 1
  AND current_timestamp - p.created < make_interval(days => $3)
ORDER BY p.created DESC;

-- name: GetPostsForHomeLatest :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       s.parent_id,
       s.name                       as space_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'post') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'post') AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id != 1
  AND current_timestamp - p.created < make_interval(days => $2)
ORDER BY p.created DESC;

-- name: GetPostsForSpacePopularByName :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       s.parent_id,
       s.name                       as space_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'post') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'post') AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on v.user_id = $3 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE s.name = $1
  AND s.parent_id = $2
  AND p.id != 1
  AND current_timestamp - p.created < make_interval(days => $4)
ORDER BY up_votes DESC;

-- name: GetPostsForSpacePopular :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       s.parent_id,
       s.name                       as space_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'post') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'post') AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on v.user_id = $2 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE space_id = $1
  AND p.id != 1
  AND current_timestamp - p.created < make_interval(days => $3)
ORDER BY up_votes DESC;

-- name: GetPostsForHomePopular :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       s.parent_id,
       s.name                       as space_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'post') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'post') AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id != 1
  AND current_timestamp - p.created
    < make_interval(days => $2)
ORDER BY up_votes DESC;

-- name: GetPostsForUser :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       s.parent_id,
       s.name                       as space_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'post') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'post') AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.poster_id = $1
  AND p.id != 1;

-- name: GetPostsForUserSubscriptionLatest :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       s.parent_id,
       s.name                       as space_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'post') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'post') AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
         join subscriptions su on su.space_id = p.space_id and su.user_id = $1 and su.is_deleted = false
WHERE p.id != 1
  AND s.id != 1
  AND current_timestamp - p.created < make_interval(days => $2)
ORDER BY p.created DESC;

-- name: GetPostsForUserSubscriptionPopular :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       s.parent_id,
       s.name                       as space_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'post') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'post') AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
         join subscriptions su on su.space_id = p.space_id and su.user_id = $1 and su.is_deleted = false
WHERE p.id != 1
  AND s.id != 1
  AND current_timestamp - p.created < make_interval(days => $2)
ORDER BY up_votes DESC;

-- name: GetPoster :one
SELECT u.*
from posts p
         join users u on p.poster_id = u.id
where p.id = $1;