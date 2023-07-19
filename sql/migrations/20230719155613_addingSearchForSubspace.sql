-- migrate:up
ALTER TABLE spaces
    ADD COLUMN ts tsvector
        GENERATED ALWAYS AS
            (setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
             setweight(to_tsvector('english', coalesce(description, '')), 'B')) STORED;
-- migrate:down

