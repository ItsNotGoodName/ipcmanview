-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_dahua_file_cursors" table
CREATE TABLE `new_dahua_file_cursors` (`device_id` integer NOT NULL, `quick_cursor` datetime NOT NULL, `full_cursor` datetime NOT NULL, `full_epoch` datetime NOT NULL, `full_complete` boolean NOT NULL AS (full_cursor <= full_epoch) STORED, `scan` boolean NOT NULL, `scan_percent` real NOT NULL, `scan_type` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- drop "dahua_file_cursors" table after copying rows
DROP TABLE `dahua_file_cursors`;
-- rename temporary table "new_dahua_file_cursors" to "dahua_file_cursors"
ALTER TABLE `new_dahua_file_cursors` RENAME TO `dahua_file_cursors`;
-- create index "dahua_file_cursors_device_id" to table: "dahua_file_cursors"
CREATE UNIQUE INDEX `dahua_file_cursors_device_id` ON `dahua_file_cursors` (`device_id`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "dahua_file_cursors_device_id" to table: "dahua_file_cursors"
DROP INDEX `dahua_file_cursors_device_id`;
-- reverse: create "new_dahua_file_cursors" table
DROP TABLE `new_dahua_file_cursors`;
