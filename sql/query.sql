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

-- name: CreateComment :one
INSERT INTO comments (parent_id, post_id, poster_id, body, content, content_type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPost :one
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true)  AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false) AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         join votes v on p.poster_id = v.user_id
WHERE p.id = $1
LIMIT 1;

-- name: GetPostsForSpaceLatest :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true)  AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false) AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on p.poster_id = v.user_id and v.user_id = $2 and p.id = v.post_or_comment_id
WHERE space_id = $1
  AND p.id != 1
  AND current_timestamp - p.created < make_interval(days => $3)
ORDER BY p.created DESC;

-- name: GetPostsForSpacePopular :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true)  AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false) AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on p.poster_id = v.user_id and v.user_id = $2 and p.id = v.post_or_comment_id
WHERE space_id = $1
  AND p.id != 1
  AND current_timestamp - p.created < make_interval(days => $3)
ORDER BY up_votes DESC;

-- name: GetPostsForUser :many
SELECT p.*,
       u.display_name,
       s.picture                    as space_picture,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true)  AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false) AS down_votes,
       v.*
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on p.poster_id = v.user_id and v.user_id = $1 and p.id = v.post_or_comment_id
WHERE p.poster_id = $1
  AND p.id != 1;

-- name: GetCommentsForPost :many
WITH RECURSIVE allCommentsForPost AS (
    -- base case starting from grandfather
    SELECT id,
           post_id,
           poster_id,
           parent_id,
           body,
           content_type,
           content,
           is_deleted,
           created,
           0 AS level
    FROM comments c
    WHERE c.post_id = $1
    UNION
    --- recursive query (note it adds to the partial table "x")
    SELECT c.id,
           c.post_id,
           c.poster_id,
           c.parent_id,
           c.body,
           c.content_type,
           c.content,
           c.is_deleted,
           c.created,
           c1.level + 1
    FROM comments c
             INNER JOIN allCommentsForPost c1
                        ON c.parent_id = c1.id
    WHERE c1.level < 3)
SELECT c.*,
       u.display_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = true)  AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = false) AS down_votes,
       v.*
FROM allCommentsForPost c
         left join votes v on c.poster_id = v.user_id and v.user_id = $2 and c.id = v.post_or_comment_id
         join users u on c.poster_id = u.id;

-- name: GetCommentsForComment :many
WITH RECURSIVE allCommentsForComment AS (
    -- base case starting from grandfather
    SELECT id,
           post_id,
           poster_id,
           parent_id,
           body,
           content_type,
           content,
           is_deleted,
           created,
           0 AS level
    FROM comments c
    WHERE c.id = $1
      AND c.id != 1
    UNION
    --- recursive query (note it adds to the partial table "x")
    SELECT c.id,
           c.post_id,
           c.poster_id,
           c.parent_id,
           c.body,
           c.content_type,
           c.content,
           c.is_deleted,
           c.created,
           c1.level + 1
    FROM comments c
             INNER JOIN allCommentsForComment c1
                        ON c.parent_id = c1.id
    WHERE c1.level < 3)
SELECT c.*,
       u.display_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = true)  AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = false) AS down_votes,
       v.*
FROM allCommentsForComment c
         join users u on c.poster_id = u.id
         left join votes v on c.poster_id = v.user_id and v.user_id = $2 and c.id = v.post_or_comment_id;
-- name: GetCommentsForUser :many
SELECT c.*,
       u.display_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = true)  AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = false) AS down_votes
FROM comments c
         join users u on c.poster_id = u.id
where u.id = $1;

-- name: GetComment :one
SELECT c.*,
       u.display_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = true)  AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = false) AS down_votes
FROM comments c
         join users u on c.poster_id = u.id
where c.id = $1;

-- name: CreateVote :one
INSERT INTO votes (user_id, post_or_comment_id, vote_type, is_up_vote)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: SetVote :exec
UPDATE votes
set is_up_vote = $1,
    is_deleted = false
where id = $2;

-- name: DeleteVote :exec
UPDATE votes
set is_deleted = true
where id = $1;

-- name: RefreshVote :exec
UPDATE votes
set is_deleted = false
where id = $1;

-- name: GetVoteForPostOrCommentForUser :one
SELECT *
FROM votes
WHERE user_id = $1
  AND post_or_comment_id = $2
  AND vote_type = $3
LIMIT 1;

-- name: GetVotesForPost :one
SELECT (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.vote_type = 'post'
          AND v.is_deleted = false
          AND v.is_up_vote = true)  AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.vote_type = 'post'
          AND v.is_deleted = false
          AND v.is_up_vote = false) AS down_votes,
       v.*
FROM posts p
         join votes v on p.id = v.post_or_comment_id
WHERE p.id = $1
  AND v.user_id = $2
  AND v.vote_type = 'post';

-- name: GetVotesForComment :one
SELECT (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.vote_type = 'comment'
          AND v.is_deleted = false
          AND v.is_up_vote = true)  AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.vote_type = 'comment'
          AND v.is_deleted = false
          AND v.is_up_vote = false) AS down_votes,
       v.*
FROM comments c
         join votes v on c.id = v.post_or_comment_id
WHERE c.id = $1
  AND v.user_id = $2
  AND v.vote_type = 'comment';
-- name: GetVoteForId :one
SELECT *
FROM votes
WHERE id = $1;