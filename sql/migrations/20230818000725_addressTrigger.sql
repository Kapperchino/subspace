-- migrate:up
CREATE OR REPLACE FUNCTION notify_addresses_update() RETURNS TRIGGER AS
$$
DECLARE
    row    RECORD;
    output TEXT;
BEGIN
    row = NEW;
    -- Forming the Output as notification. You can choose you own notification.
    output = json_build_object('from_id', row.from_user_id, 'to_id', row.to_user_id, 'post_id', row.post_id);

    -- Calling the pg_notify for my_table_update event with output as payload

    PERFORM pg_notify('addresses_update', output);

    -- Returning null because it is an after trigger.
    RETURN NULL;
END ;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_user_addresses_update
    AFTER INSERT OR UPDATE
    ON user_addresses
    FOR EACH ROW
EXECUTE PROCEDURE notify_addresses_update();
-- migrate:down

