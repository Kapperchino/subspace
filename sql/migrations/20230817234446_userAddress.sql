-- migrate:up
CREATE TABLE user_addresses
(
    id           BIGSERIAL PRIMARY KEY,
    from_user_id BIGINT REFERENCES users (id),
    to_user_id   BIGINT REFERENCES users (id),
    post_id      BIGINT REFERENCES posts (id)
);
-- migrate:down

