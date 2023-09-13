-- migrate:up
ALTER TABLE user_addresses
    ADD COLUMN comment_id BIGINT REFERENCES comments (id)
-- migrate:down

