CREATE TABLE IF NOT EXISTS data_table (
    id bigserial not null constraint data_table_pk primary key,
    name varchar(255) not null,
    schema jsonb not null,
    account_uuid uuid not null,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);

CREATE TABLE IF NOT EXISTS source_table_map (
    id bigserial not null constraint source_table_map_pk primary key,
    table_id int not null,
    source_id int not null,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
    );

CREATE TABLE IF NOT EXISTS connection (
    id bigserial not null constraint source_connection_pk primary key,
    name varchar(255) not null,
    type varchar(255) not null,
    configurations jsonb not null,
    account_uuid uuid not null,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
    );