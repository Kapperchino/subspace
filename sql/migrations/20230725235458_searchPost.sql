-- migrate:up
ALTER TABLE posts
    ADD COLUMN ts tsvector
        GENERATED ALWAYS AS
            (setweight(to_tsvector('english', coalesce(body, '')), 'A') ||
             setweight(to_tsvector('english', coalesce(topic, '')), 'B')) STORED;
-- migrate:down
