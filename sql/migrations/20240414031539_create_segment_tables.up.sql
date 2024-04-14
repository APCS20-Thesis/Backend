ALTER TABLE data_action RENAME COLUMN source_table_map_id TO object_id;
ALTER TABLE data_action ADD COLUMN object_type varchar(255);

CREATE TABLE IF NOT EXISTS master_segment (
    id bigserial not null constraint master_segment_pk primary key,
    name varchar(255) not null,
    description varchar(255),
    status varchar(255) not null ,
    account_uuid uuid not null,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);

CREATE TABLE IF NOT EXISTS segment (
    id bigserial not null constraint segment_pk primary key,
    master_segment_id bigint not null,
    name varchar(255) not null,
    description varchar(255),
    condition jsonb,
    account_uuid uuid not null,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);

CREATE TABLE IF NOT EXISTS audience_table (
    id bigserial not null constraint audience_table_pk primary key,
    master_segment_id bigint not null,
    name varchar(255) not null,
    schema jsonb,
    build_configuration jsonb,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);

CREATE TABLE IF NOT EXISTS behavior_table (
    id bigserial not null constraint behavior_table_pk primary key,
    master_segment_id bigint not null,
    data_table_id bigint not null,
    name varchar(255) not null,
    schema jsonb,
    audience_foreign_key varchar(255),
    join_key varchar(255),
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);