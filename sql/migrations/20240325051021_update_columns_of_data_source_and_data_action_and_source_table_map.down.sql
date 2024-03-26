ALTER TABLE data_source RENAME COLUMN configurations TO configuration;
ALTER TABLE data_source ADD COLUMN mapping_options jsonb;
ALTER TABLE data_source DROP COLUMN  status;

ALTER TABLE data_table DROP COLUMN  status;
ALTER TABLE data_table ALTER COLUMN schema SET NOT NULL;

ALTER TABLE source_table_map DROP COLUMN mapping_options;
ALTER TABLE source_table_map DROP COLUMN table_name_in_source;

ALTER TABLE data_action DROP COLUMN source_table_map_id;