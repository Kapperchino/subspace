-- migrate:up
CREATE OR REPLACE VIEW posts_view AS
SELECT p.*,
       u.display_name,
       user_pic.url                 as user_pic_url,
       user_pic.width               as user_pic_width,
       user_pic.height              as user_pic_height,
       user_pic.id                  as user_pic_id,
       space_small_pic.url          as space_small_pic_url,
       space_small_pic.width        as space_small_pic_width,
       space_small_pic.height       as space_small_pic_height,
       space_small_pic.id           as space_small_pic_id,
       s.parent_id,
       s.name                       as space_name,
       (SELECT COUNT(*)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = true
          AND v.vote_type = 'post') AS up_votes,
       (SELECT COUNT(*)
        FROM votes v
        WHERE v.post_or_comment_id = p.id
          AND v.is_deleted = false
          AND v.is_up_vote = false
          AND v.vote_type = 'post') AS down_votes,
       (SELECT COUNT(*)
        FROM comments c
        WHERE c.post_id = p.id
          AND c.is_deleted = false) AS comment_count
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join pictures user_pic on u.picture_id = user_pic.id
         left join pictures space_small_pic on s.small_picture_id = space_small_pic.id;

CREATE OR REPLACE VIEW spaces_view AS
SELECT s.*,
       small_picture.url              as space_small_pic_url,
       small_picture.width            as space_small_pic_width,
       small_picture.height           as space_small_pic_height,
       small_picture.id               as space_small_picture_id,
       background_picture.url         as space_background_picture_url,
       background_picture.width       as space_background_picture_width,
       background_picture.height      as space_background_picture_height,
       background_picture.id          as space_background_picture_id,
       (SELECT COUNT(*)
        FROM subscriptions sub
        WHERE sub.space_id = s.id
          AND sub.is_deleted = false) as sub_count
FROM spaces s
         left join pictures small_picture on s.small_picture_id = small_picture.id
         left join pictures background_picture on s.background_picture_id = background_picture.id;
-- migrate:down

