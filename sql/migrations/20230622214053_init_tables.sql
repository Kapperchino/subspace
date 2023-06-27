-- migrate:up
CREATE TABLE spaces
(
    id          BIGSERIAL PRIMARY KEY,
    parent_id   BIGSERIAL NOT NULL,
    name        TEXT      NOT NULL,
    description TEXT,
    picture     TEXT,
    is_deleted  BOOLEAN DEFAULT false,
    UNIQUE (parent_id, name)
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
    display_name TEXT NOT NULL UNIQUE,
    bio          TEXT,
    is_deleted   BOOLEAN DEFAULT false
);

INSERT INTO users (password, email, display_name)
VALUES ('joe', 'joebiden123123@gmail.com', 'root');

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
    is_deleted   BOOLEAN DEFAULT false
);

INSERT INTO posts (space_id, poster_id, topic)
VALUES (1, 1, 'rootPost');

CREATE TABLE comments
(
    id           BIGSERIAL PRIMARY KEY,
    post_id      BIGINT REFERENCES posts (id),
    poster_id    BIGINT REFERENCES users (id),
    parent_id    BIGINT,
    body         TEXT NOT NULL,
    content_type content_type,
    content      TEXT,
    is_deleted   BOOLEAN DEFAULT false
);

INSERT INTO comments(post_id, poster_id, parent_id, body)
VALUES (1, 1, 1, 'root');

ALTER TABLE comments
    ADD CONSTRAINT comments_self_ref foreign key (parent_id) references comments (id);

CREATE TYPE vote_type AS ENUM ('post', 'comment');

CREATE TABLE votes
(
    id                 BIGSERIAL PRIMARY KEY,
    is_up_vote         BOOLEAN DEFAULT true,
    user_id            BIGINT REFERENCES users (id) NOT NULL,
    post_or_comment_id BIGINT                       NOT NULL,
    vote_type          vote_type                    NOT NULL,
    is_deleted         BOOLEAN DEFAULT false
);

CREATE TABLE subscriptions
(
    id       BIGSERIAL PRIMARY KEY,
    user_id  BIGINT REFERENCES users (id) NOT NULL,
    space_id BIGINT REFERENCES spaces (id)
);
-- migrate:down

