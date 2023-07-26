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
         join tags_relations t on t.post_or_comment_id = p.id
         join tags t1 on t.tag_id = t1.id
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on p.poster_id = v.user_id and v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id != 1
  AND t1.name = $2
  AND current_timestamp - p.created < make_interval(days => $3)
ORDER BY up_votes DESC;

-- name: GetPostsWithTagsLatest :many
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
         join tags_relations t on t.post_or_comment_id = p.id
         join tags t1 on t.tag_id = t1.id
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join votes v on p.poster_id = v.user_id and v.user_id = $1 and p.id = v.post_or_comment_id and
                              v.vote_type = 'post'
WHERE p.id != 1
  AND t1.name = $2
  AND current_timestamp - p.created < make_interval(days => $3)
ORDER BY p.created DESC;