-- migrate:up
CREATE TABLE mentions
(
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT REFERENCES users (id),
    user_mentioned_id BIGINT REFERENCES users (id)
);
-- migrate:down

