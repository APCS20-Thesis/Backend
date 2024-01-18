CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS  account (
    id bigserial not null constraint user_pk primary key,
    uuid uuid DEFAULT uuid_generate_v4 (),
    username varchar(255) unique not null,
    password varchar(128) not null,
    first_name varchar(255),
    last_name varchar(255),
    email varchar UNIQUE NOT NULL,
    created_at    timestamptz default (NOW() AT TIME ZONE 'UTC') not null,
    updated_at    timestamptz default (NOW() AT TIME ZONE 'UTC') not null
);