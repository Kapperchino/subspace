-- migrate:up
ALTER TABLE subscriptions
    ADD CONSTRAINT subscriptions_unique UNIQUE (user_id, space_id);
-- migrate:down

