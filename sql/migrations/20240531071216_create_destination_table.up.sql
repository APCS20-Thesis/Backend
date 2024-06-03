CREATE TABLE IF NOT EXISTS data_destination (
    id bigserial not null constraint data_destination_pk primary key,
    name varchar(255) not null,
    account_uuid uuid not null,
    type varchar(255) not null,
    configurations jsonb not null,
    connection_id int,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);

ALTER TABLE data_source ADD COLUMN IF NOT EXISTS connection_id int;