-- migrate:up
CREATE EXTENSION pg_trgm;
CREATE INDEX space_name_idx ON spaces USING GIST (name gist_trgm_ops);
CREATE INDEX user_ad_idx on users USING GIST (address gist_trgm_ops);
-- migrate:down

