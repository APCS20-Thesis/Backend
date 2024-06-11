CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS set_timestamp_account on "account";
CREATE TRIGGER set_timestamp_account
    BEFORE UPDATE ON "account"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


DROP TRIGGER IF EXISTS set_timestamp_audience_table on "audience_table";
CREATE TRIGGER set_timestamp_audience_table
    BEFORE UPDATE ON "audience_table"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


DROP TRIGGER IF EXISTS set_timestamp_behavior_table on "behavior_table";
CREATE TRIGGER set_timestamp_behavior_table
    BEFORE UPDATE ON "behavior_table"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


DROP TRIGGER IF EXISTS set_timestamp_connection on "connection";
CREATE TRIGGER set_timestamp_connection
    BEFORE UPDATE ON "connection"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


DROP TRIGGER IF EXISTS set_timestamp_data_action on "data_action";
CREATE TRIGGER set_timestamp_data_action
    BEFORE UPDATE ON "data_action"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


DROP TRIGGER IF EXISTS set_timestamp_data_action_run on "data_action_run";
CREATE TRIGGER set_timestamp_data_action_run
    BEFORE UPDATE ON "data_action_run"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


DROP TRIGGER IF EXISTS set_timestamp_data_destination on "data_destination";
CREATE TRIGGER set_timestamp_data_destination
    BEFORE UPDATE ON "data_destination"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


DROP TRIGGER IF EXISTS set_timestamp_data_source on "data_source";
CREATE TRIGGER set_timestamp_data_source
    BEFORE UPDATE ON "data_source"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


DROP TRIGGER IF EXISTS set_timestamp_data_table on "data_table";
CREATE TRIGGER set_timestamp_data_table
    BEFORE UPDATE ON "data_table"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


DROP TRIGGER IF EXISTS set_timestamp_master_segment on "master_segment";
CREATE TRIGGER set_timestamp_master_segment
    BEFORE UPDATE ON "master_segment"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


DROP TRIGGER IF EXISTS set_timestamp_segment on "segment";
CREATE TRIGGER set_timestamp_segment
    BEFORE UPDATE ON "segment"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


DROP TRIGGER IF EXISTS set_timestamp_source_table_map on "source_table_map";
CREATE TRIGGER set_timestamp_source_table_map
    BEFORE UPDATE ON "source_table_map"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


DROP TRIGGER IF EXISTS set_timestamp_dest_table_map on "dest_table_map";
CREATE TRIGGER set_timestamp_dest_table_map
    BEFORE UPDATE ON "dest_table_map"
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
