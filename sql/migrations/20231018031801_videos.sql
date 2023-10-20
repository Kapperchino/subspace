-- migrate:up
CREATE TYPE process_state AS ENUM ('error','created', 'done', 'ongoing');
CREATE TABLE videos
(
    id            BIGSERIAL PRIMARY KEY,
    url           TEXT NOT NULL,
    thumbnail     TEXT,
    process_state process_state,
    stream_url    TEXT,
    duration      double precision,
    width         INT,
    height        INT
);

CREATE TABLE video_relations
(
    id         BIGSERIAL PRIMARY KEY,
    video_id   BIGINT REFERENCES videos (id),
    post_id    BIGINT,
    comment_id BIGINT
);
-- migrate:down

