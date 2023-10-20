-- name: CreateVideo :one
INSERT INTO videos (url, process_state)
VALUES ($1, $2)
RETURNING *;

-- name: CreateVideoRelation :one
INSERT INTO video_relations (video_id, post_id, comment_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateVideo :exec
UPDATE videos
set stream_url    = $1,
    process_state = $2,
    thumbnail     = $3,
    duration      = $4,
    width         = $5,
    height        = $6
where id = $7;

-- name: GetVideosForPost :many
SELECT v.*
FROM videos v
         join video_relations vr on v.id = vr.video_id
where vr.post_id = $1;

-- name: GetVideo :one
SELECT v.*
FROM videos v
where v.id = $1;

-- name: GetVideoForComment :many
SELECT v.*
FROM videos v
         join video_relations vr on v.id = vr.video_id
where vr.post_id = $1;