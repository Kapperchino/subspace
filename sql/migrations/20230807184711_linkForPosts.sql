-- migrate:up
ALTER TABLE pictures
    DROP COLUMN post_id,
    DROP COLUMN comment_id;

ALTER TABLE posts
    ADD COLUMN link TEXT;

-- migrate:down

