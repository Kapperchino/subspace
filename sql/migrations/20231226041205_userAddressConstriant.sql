-- migrate:up
ALTER TABLE users
    ADD CONSTRAINT user_address_unique UNIQUE(address);
-- migrate:down

