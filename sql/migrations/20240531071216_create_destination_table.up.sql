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

CREATE TABLE IF NOT EXISTS dest_table_map (
  id bigserial not null constraint dest_table_map_pk primary key,
  table_id int not null,
  destination_id int not null,
  mapping_options jsonb,
  dest_table_name varchar(255),
  created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
  updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);

CREATE TABLE IF NOT EXISTS dest_segment_map (
  id bigserial not null constraint dest_segment_map_pk primary key,
  segment_id int not null,
  destination_id int not null,
  mapping_options jsonb,
  created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
  updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);

CREATE TABLE IF NOT EXISTS dest_ms_segment_map (
  id bigserial not null constraint dest_ms_segment_map_pk primary key,
  master_segment_id int not null,
  destination_id int not null,
  mapping_options jsonb,
  created_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null,
  updated_at timestamptz default (NOW () AT TIME ZONE 'UTC') not null
);
