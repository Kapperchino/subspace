-- migrate:up
ALTER TABLE posts ALTER COLUMN topic DROP NOT NULL;
-- migrate:down

