DROP TRIGGER IF EXISTS set_timestamp_setting on "setting";

DROP TABLE IF EXISTS setting;

ALTER TABLE account DROP COLUMN phone;
ALTER TABLE account DROP COLUMN country;
ALTER TABLE account DROP COLUMN company;
ALTER TABLE account DROP COLUMN position;