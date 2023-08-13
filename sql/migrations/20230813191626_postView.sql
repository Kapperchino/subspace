-- migrate:up
CREATE VIEW posts_view AS
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
          AND v.vote_type = 'post') AS down_votes
FROM posts p
         join users u on p.poster_id = u.id
         join spaces s on s.id = p.space_id
         left join pictures user_pic on u.picture_id = user_pic.id
         left join pictures space_small_pic on s.small_picture_id = space_small_pic.id;
-- migrate:down

