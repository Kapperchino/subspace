-- migrate:up
CREATE TABLE pictures
(
    id         BIGSERIAL PRIMARY KEY,
    post_id    BIGINT,
    comment_id BIGINT,
    url        TEXT   NOT NULL,
    width      BIGINT NOT NULL,
    height     BIGINT NOT NULL
);

CREATE TABLE pictureRelations
(
    id         BIGSERIAL PRIMARY KEY,
    picture_id BIGINT REFERENCES pictures (id),
    post_id    BIGINT,
    comment_id BIGINT
);

ALTER TABLE posts
    DROP COLUMN content;
ALTER TABLE comments
    DROP COLUMN content;
ALTER TABLE users
    DROP COLUMN picture;
ALTER TABLE spaces
    DROP COLUMN picture;

ALTER TABLE users
    ADD COLUMN picture_id BIGINT;
ALTER TABLE spaces
    ADD COLUMN small_picture_id      BIGINT,
    ADD COLUMN background_picture_id BIGINT;
-- migrate:down

