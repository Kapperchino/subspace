-- migrate:up
ALTER TABLE users
    ADD COLUMN address text;

ALTER TABLE users
    ADD COLUMN ts tsvector
        GENERATED ALWAYS AS
            (setweight(to_tsvector('english', coalesce(display_name, '')), 'A') ||
             setweight(to_tsvector('english', coalesce(address, '')), 'B')) STORED;
-- migrate:down

