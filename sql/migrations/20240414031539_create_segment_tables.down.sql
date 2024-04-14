ALTER TABLE data_action RENAME COLUMN object_id TO source_table_map_id;
ALTER TABLE data_action DROP COLUMN target_table;

DROP TABLE IF EXISTS audience_table;
DROP TABLE IF EXISTS behavior_table;
DROP TABLE IF EXISTS segment;
DROP TABLE IF EXISTS master_segment;
