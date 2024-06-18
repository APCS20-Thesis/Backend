DROP TABLE IF EXISTS data_destination;
DROP TABLE IF EXISTS dest_table_map;
DROP TABLE IF EXISTS dest_segment_map;
DROP TABLE IF EXISTS dest_ms_segment_map;

ALTER TABLE data_source DROP COLUMN IF EXISTS connection_id;

