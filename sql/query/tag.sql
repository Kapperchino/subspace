-- name: CreateTag :one
INSERT INTO tags (name)
VALUES ($1)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetTag :one
SELECT *
FROM tags
where name = $1;

-- name: CreateTagRelation :one
INSERT INTO tags_relations (tag_id, post_or_comment_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetPostsWithTagsPopular :many
SELECT sqlc.embed(p), v.*
FROM posts_view p
         join tags_relations t on t.post_or_comment_id = p.id
         join tags t1 on t.tag_id = t1.id
         left join votes v on p.poster_id = v.user_id and v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id != 1
  AND t1.name = $2
  AND current_timestamp - p.created < make_interval(days => $3)
ORDER BY up_votes DESC;

-- name: GetPostsWithTagsLatest :many
SELECT sqlc.embed(p), v.*
FROM posts_view p
         join tags_relations t on t.post_or_comment_id = p.id
         join tags t1 on t.tag_id = t1.id
         left join votes v on p.poster_id = v.user_id and v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id != 1
  AND t1.name = $2
  AND current_timestamp - p.created < make_interval(days => $3)
ORDER BY p.created DESC;