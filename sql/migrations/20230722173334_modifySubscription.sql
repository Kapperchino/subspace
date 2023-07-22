-- migrate:up
ALTER TABLE subscriptions
    ADD COLUMN is_deleted BOOLEAN DEFAULT false;
-- migrate:down

