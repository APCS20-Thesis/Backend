ALTER TABLE source_table_map RENAME COLUMN table_name_in_source TO source_table_name;

CREATE TABLE IF NOT EXISTS file_export_record (
    id bigserial not null constraint file_export_record_pk primary key,
    data_table_id int,
    format varchar(255),
    account_uuid uuid not null,
    data_action_id int,
    data_action_run_id int,
    status varchar(255),
    s3_key varchar(255),
    download_url varchar(255),
    expiration_time timestamptz,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);