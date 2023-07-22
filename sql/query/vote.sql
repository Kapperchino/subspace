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
         join votes v on p.id = v.post_or_comment_id and v.vote_type = 'post'
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
         join votes v on c.id = v.post_or_comment_id and v.vote_type = 'comment'
WHERE c.id = $1
  AND v.user_id = $2
  AND v.vote_type = 'comment';
-- name: GetVoteForId :one
SELECT *
FROM votes
WHERE id = $1;