-- migrate:up
ALTER TABLE devices
    ADD COLUMN device_info TEXT;
-- migrate:down

