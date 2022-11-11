CREATE table test (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE OR REPLACE FUNCTION notify_test_update() RETURNS TRIGGER AS $$
DECLARE
    row RECORD;
BEGIN
    IF (TG_OP = 'DELETE') THEN
        row = OLD;
    ELSE
        row = NEW;
    END IF;

    PERFORM pg_notify(
        'test_table_update',
        (
            SELECT row_to_json(payload)
            FROM (
                SELECT row_to_json(row) AS "data",
                       TG_OP AS "operation") payload)::TEXT

    );
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_test_update
    AFTER INSERT OR UPDATE OR DELETE
    ON test
    FOR EACH ROW
EXECUTE PROCEDURE notify_test_update();