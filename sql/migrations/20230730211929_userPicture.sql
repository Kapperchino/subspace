-- migrate:up
ALTER TABLE users
    ADD COLUMN picture TEXT;
-- migrate:down

