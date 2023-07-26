-- migrate:up
CREATE TABLE tags
(
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    UNIQUE (name)
);

CREATE TABLE tags_relations
(
    id                 BIGSERIAL PRIMARY KEY,
    tag_id             BIGINT REFERENCES tags (id),
    post_or_comment_id BIGINT NOT NULL
);
-- migrate:down

