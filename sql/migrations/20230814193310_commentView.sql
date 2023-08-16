-- migrate:up
CREATE VIEW comments_view AS
SELECT c.*,
       user_pic.url                    as user_pic_url,
       user_pic.width                  as user_pic_width,
       user_pic.height                 as user_pic_height,
       user_pic.id                     as user_pic_id,
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
         left join pictures user_pic on u.picture_id = user_pic.id;

ALTER TABLE comments
    ADD COLUMN link TEXT;
-- migrate:down

