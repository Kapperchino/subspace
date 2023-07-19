-- migrate:up
UPDATE spaces SET name = 'SubSpace' WHERE id = 1;
-- migrate:down

