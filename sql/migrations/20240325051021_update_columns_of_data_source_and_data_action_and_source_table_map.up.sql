ALTER TABLE data_source RENAME COLUMN configuration TO configurations;
ALTER TABLE data_source DROP COLUMN mapping_options;
ALTER TABLE data_source ADD COLUMN  status varchar(255) not null;

ALTER TABLE data_table add column status varchar(255) not null;
ALTER TABLE data_table ALTER COLUMN schema DROP NOT NULL;

ALTER TABLE source_table_map ADD COLUMN mapping_options jsonb;
ALTER TABLE source_table_map ADD COLUMN table_name_in_source varchar(255);

ALTER TABLE data_action ADD COLUMN source_table_map_id int not null;