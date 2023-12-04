-- name: CreateTag :one
INSERT INTO tags (name)
VALUES ($1)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetTag :one
SELECT *
FROM tags
where name = $1;

-- name: CreateTagRelationForPost :one
INSERT INTO tags_relations (tag_id, post_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetPostsWithTagsPopular :many
SELECT sqlc.embed(p), v.*
FROM posts_view p
         join tags_relations t on t.post_id = p.id
         join tags t1 on t.tag_id = t1.id
         left join votes v on p.poster_id = v.user_id and v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id != 1
  AND t1.name = $2
  AND current_timestamp - p.created < make_interval(days => $3)
ORDER BY up_votes DESC
LIMIT 10 OFFSET $4;

-- name: GetPostsWithTagsLatest :many
SELECT sqlc.embed(p), v.*
FROM posts_view p
         join tags_relations t on t.post_id = p.id
         join tags t1 on t.tag_id = t1.id
         left join votes v on p.poster_id = v.user_id and v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id != 1
  AND t1.name = $2
  AND current_timestamp - p.created < make_interval(days => $3)
ORDER BY p.created DESC
LIMIT 10 OFFSET $4;

-- name: GetPopularTags :many
select t.name, COUNT(DISTINCT tr.post_id) AS count
from tags_relations tr
         join posts_view pv on pv.id = tr.post_id
         join tags t on t.id = tr.tag_id
WHERE pv.id != 1
  AND current_timestamp - pv.created < make_interval(days => $1)
GROUP BY tr.tag_id, t.name
ORDER BY COUNT(tr.tag_id) DESC;

-- name: SearchPrefixTags :many
select sqlc.embed(t), similarity(t.name, $1) as sml
from tags t
order by sml desc
limit 5;