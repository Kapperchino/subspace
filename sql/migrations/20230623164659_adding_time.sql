-- migrate:up
ALTER TABLE spaces ADD COLUMN created timestamp default current_timestamp;
ALTER TABLE posts ADD COLUMN created timestamp default current_timestamp;
ALTER TABLE comments ADD COLUMN created timestamp default current_timestamp;
ALTER TABLE users ADD COLUMN created timestamp default current_timestamp;

ALTER TABLE spaces ADD UNIQUE(name);
ALTER TABLE users ADD UNIQUE(display_name);
-- migrate:down

