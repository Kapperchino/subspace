-- migrate:up
CREATE OR REPLACE FUNCTION notify_comments_update() RETURNS TRIGGER AS
$$
DECLARE
    row    RECORD;
    output TEXT;
BEGIN
    row = NEW;
    -- Forming the Output as notification. You can choose you own notification.
    output = json_build_object('id', row.id, 'parent_id', row.parent_id, 'post_id', row.post_id);

    -- Calling the pg_notify for my_table_update event with output as payload

    PERFORM pg_notify('comments_update', output);

    -- Returning null because it is an after trigger.
    RETURN NULL;
END ;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_comment_update
    AFTER INSERT OR UPDATE
    ON comments
    FOR EACH ROW
EXECUTE PROCEDURE notify_comments_update();
-- migrate:down

