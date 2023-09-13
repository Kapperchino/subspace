-- migrate:up
ALTER TABLE tags_relations
    DROP COLUMN post_or_comment_id;

ALTER TABLE tags_relations
    ADD COLUMN post_id BIGINT references posts (id);

ALTER TABLE tags_relations
    ADD COLUMN comment_id BIGINT references comments (id);
-- migrate:down

