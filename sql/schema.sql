CREATE TABLE spaces
(
    id          BIGSERIAL PRIMARY KEY,
    parent_id   BIGSERIAL REFERENCES spaces (id),
    name        text NOT NULL,
    description text
);

CREATE TABLE users
(
    id           BIGSERIAL PRIMARY KEY,
    email        text NOT NULL,
    display_name text NOT NULL,
    bio          text
);

CREATE TABLE posts
(
    id         BIGSERIAL PRIMARY KEY,
    space_id   BIGINT REFERENCES spaces (id),
    poster_id  BIGINT REFERENCES users (id),
    topic      text NOT NULL,
    content    text,
    up_votes   INTEGER,
    down_votes INTEGER
);

CREATE TABLE comments
(
    id         BIGSERIAL PRIMARY KEY,
    post_id    BIGINT REFERENCES posts (id),
    poster_id  BIGINT REFERENCES users (id),
    parent_id  BIGINT REFERENCES comments (id),
    content    text,
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

