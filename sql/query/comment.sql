-- name: CreateComment :one
INSERT INTO comments (parent_id, post_id, poster_id, body, content_type)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetCommentsForPost :many
WITH RECURSIVE allCommentsForPost AS (SELECT c.id, 0 AS level
                                      FROM comments c
                                      WHERE c.post_id = $1
                                      UNION
                                      SELECT c.id,
                                             c1.level + 1
                                      FROM comments c
                                               INNER JOIN allCommentsForPost c1
                                                          ON c.parent_id = c1.id
                                      WHERE c1.level < 3)
SELECT sqlc.embed(cv), v.*
from comments_view cv
         left join votes v on v.user_id = $2 and cv.id = v.post_or_comment_id and
                              v.vote_type = 'comment'
where cv.id IN (SELECT id from allCommentsForPost);

-- name: GetCommentsForComment :many
WITH RECURSIVE allCommentsForComment AS (SELECT c.id,
                                                0 AS level
                                         FROM comments c
                                         WHERE c.id = $1
                                           AND c.id != 1
                                         UNION
                                         SELECT c.id,
                                                c1.level + 1
                                         FROM comments c
                                                  INNER JOIN allCommentsForComment c1
                                                             ON c.parent_id = c1.id
                                         WHERE c1.level < 3)
SELECT sqlc.embed(cv), v.*
from comments_view cv
         left join votes v on v.user_id = $2 and cv.id = v.post_or_comment_id and
                              v.vote_type = 'comment'
where cv.id IN (allCommentsForComment);

-- name: GetCommentsForUser :many
SELECT sqlc.embed(c)
FROM comments_view c
         left join votes v on v.user_id = $2 and c.id = v.post_or_comment_id and
                              v.vote_type = 'comment'
where c.poster_id = $1;

-- name: GetComment :one
SELECT sqlc.embed(c)
FROM comments_view c
where c.id = $1;

-- name: GetCommenter :one
SELECT u.*
from comments c
         join users u on c.poster_id = u.id
where c.id = $1;