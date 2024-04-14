ALTER TABLE source_table_map RENAME COLUMN source_table_name TO table_name_in_source;

DROP TABLE IF EXISTS file_export_record;