-- name: CreateComment :one
INSERT INTO comments (parent_id, post_id, poster_id, body, content_type)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCommentsForPost :many
WITH RECURSIVE allCommentsForPost AS (
    -- base case starting from grandfather
    SELECT id,
           post_id,
           poster_id,
           parent_id,
           body,
           content_type,
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
           c.is_deleted,
           c.created,
           c1.level + 1
    FROM comments c
             INNER JOIN allCommentsForPost c1
                        ON c.parent_id = c1.id
    WHERE c1.level < 3)
SELECT sqlc.embed(c),
       
       v.*,
       u.display_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'comment') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'comment') AS down_votes
FROM allCommentsForPost c
         left join votes v on v.user_id = $2 and c.id = v.post_or_comment_id and
                              v.vote_type = 'comment'
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
           c.is_deleted,
           c.created,
           c1.level + 1
    FROM comments c
             INNER JOIN allCommentsForComment c1
                        ON c.parent_id = c1.id
    WHERE c1.level < 3)
SELECT sqlc.embed(c),
       
       v.*,
       u.display_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'comment') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'comment') AS down_votes
FROM allCommentsForComment c
         join users u on c.poster_id = u.id
         left join votes v on v.user_id = $2 and c.id = v.post_or_comment_id and
                              v.vote_type = 'comment';
-- name: GetCommentsForUser :many
SELECT sqlc.embed(c),
       
       u.display_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'comment') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'comment') AS down_votes
FROM comments c
         join users u on c.poster_id = u.id
where u.id = $1;
-- name: GetComment :one
SELECT sqlc.embed(c),
       
       u.display_name,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'comment') AS up_votes,
       (SELECT COUNT(id)
        FROM votes v
        WHERE v.post_or_comment_id = c.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'comment') AS down_votes
FROM comments c
         join users u on c.poster_id = u.id
where c.id = $1;

-- name: GetCommenter :one
SELECT u.*
from comments c
         join users u on c.poster_id = u.id
where c.id = $1;