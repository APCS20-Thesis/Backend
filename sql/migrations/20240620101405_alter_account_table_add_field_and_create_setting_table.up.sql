ALTER TABLE account ADD COLUMN phone varchar(20);
ALTER TABLE account ADD COLUMN country varchar(2);
ALTER TABLE account ADD COLUMN company varchar(50);
ALTER TABLE account ADD COLUMN position varchar(50);

CREATE TABLE IF NOT EXISTS setting (
    id bigserial not null constraint setting_pk primary key,
    notify_create_source boolean default true,
    notify_create_destination boolean default true,
    notify_create_master_segment boolean default true,
    notify_create_segment boolean default true,
    account_uuid uuid not null,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);

DROP TRIGGER IF EXISTS set_timestamp_setting on "setting";
CREATE TRIGGER set_timestamp_setting
    BEFORE UPDATE ON "setting"
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();


INSERT INTO setting(account_uuid)
SELECT uuid  FROM account WHERE uuid NOT IN (SELECT account_uuid FROM setting);