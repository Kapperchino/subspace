-- migrate:up
CREATE TABLE spaces
(
    id          BIGSERIAL PRIMARY KEY,
    parent_id   BIGSERIAL NOT NULL,
    name        TEXT      NOT NULL,
    description TEXT
);

INSERT INTO spaces (name, description, parent_id)
VALUES ('root', 'base space', 1);

ALTER TABLE spaces
    ADD CONSTRAINT spaces_self_ref foreign key (parent_id) references spaces (id);

CREATE TABLE users
(
    id           BIGSERIAL PRIMARY KEY,
    password     TEXT NOT NULL,
    email        TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    bio          TEXT
);

CREATE TYPE content_type AS ENUM ('video', 'picture', 'text');

CREATE TABLE posts
(
    id           BIGSERIAL PRIMARY KEY,
    space_id     BIGINT REFERENCES spaces (id),
    poster_id    BIGINT REFERENCES users (id),
    topic        TEXT NOT NULL,
    body         TEXT,
    content_type content_type,
    content      TEXT,
    up_votes     INTEGER,
    down_votes   INTEGER
);

CREATE TABLE comments
(
    id         BIGSERIAL PRIMARY KEY,
    post_id    BIGINT REFERENCES posts (id),
    poster_id  BIGINT REFERENCES users (id),
    parent_id  BIGINT REFERENCES comments (id),
    content    TEXT,
    up_votes   INTEGER,
    down_votes INTEGER
);

CREATE TABLE likes
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users (id) NOT NULL,
    post_id    BIGINT REFERENCES posts (id),
    comment_id BIGINT REFERENCES comments (id),
    CONSTRAINT check_single_source CHECK (num_nonnulls(post_id, comment_id) = 1)
);

CREATE TABLE dislikes
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users (id) NOT NULL,
    post_id    BIGINT REFERENCES posts (id),
    comment_id BIGINT REFERENCES comments (id),
    CONSTRAINT check_single_source CHECK (num_nonnulls(post_id, comment_id) = 1)
);

CREATE TABLE subscriptions
(
    id       BIGSERIAL PRIMARY KEY,
    user_id  BIGINT REFERENCES users (id) NOT NULL,
    space_id BIGINT REFERENCES spaces (id)
);
-- migrate:down

