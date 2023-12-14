-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_dahua_file_scan_locks" table
CREATE TABLE `new_dahua_file_scan_locks` (`camera_id` integer NOT NULL, `touched_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`camera_id`) REFERENCES `dahua_cameras` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- copy rows from old table "dahua_file_scan_locks" to new temporary table "new_dahua_file_scan_locks"
INSERT INTO `new_dahua_file_scan_locks` (`camera_id`, `touched_at`) SELECT `camera_id`, `touched_at` FROM `dahua_file_scan_locks`;
-- drop "dahua_file_scan_locks" table after copying rows
DROP TABLE `dahua_file_scan_locks`;
-- rename temporary table "new_dahua_file_scan_locks" to "dahua_file_scan_locks"
ALTER TABLE `new_dahua_file_scan_locks` RENAME TO `dahua_file_scan_locks`;
-- create index "dahua_file_scan_locks_camera_id" to table: "dahua_file_scan_locks"
CREATE UNIQUE INDEX `dahua_file_scan_locks_camera_id` ON `dahua_file_scan_locks` (`camera_id`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "dahua_file_scan_locks_camera_id" to table: "dahua_file_scan_locks"
DROP INDEX `dahua_file_scan_locks_camera_id`;
-- reverse: create "new_dahua_file_scan_locks" table
DROP TABLE `new_dahua_file_scan_locks`;
