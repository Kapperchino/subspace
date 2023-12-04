-- migrate:up
CREATE INDEX tag_name_idx ON tags USING GIST (name gist_trgm_ops);
-- migrate:down

