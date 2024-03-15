CREATE TABLE IF NOT EXISTS source_connection (
    id bigserial not null constraint source_connection_pk primary key,
    name varchar(255) not null,
    type varchar(255) not null,
    configurations jsonb not null,
    account_uuid uuid not null,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
    );