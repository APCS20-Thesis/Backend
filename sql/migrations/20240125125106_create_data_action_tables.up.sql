CREATE TABLE IF NOT EXISTS data_source (
    id bigserial not null constraint data_source_pk primary key,
    name varchar(255) not null,
    description varchar(255),
    type varchar(255) not null,
    configuration jsonb not null,
    mapping_options jsonb,
    delta_table_name varchar(255) not null,
    account_uuid uuid not null,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);

CREATE TABLE IF NOT EXISTS data_action (
    id bigserial not null constraint data_action_pk primary key,
    action_type varchar(255) not null,
    payload jsonb,
    status varchar(255) not null ,
    run_count int not null,
    schedule varchar(25),
    dag_id varchar(255),
    account_uuid uuid not null,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);

CREATE TABLE IF NOT EXISTS data_action_run (
    id bigserial not null constraint data_action_run_pk primary key,
    action_id bigint not null,
    run_id int not null,
    status varchar(255) not null,
    error varchar(255),
    account_uuid uuid not null,
    created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
    updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);

