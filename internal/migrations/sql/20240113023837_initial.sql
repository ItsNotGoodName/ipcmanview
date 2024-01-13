-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- drop "dahua_file_scan_locks" table
DROP TABLE `dahua_file_scan_locks`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: drop "dahua_file_scan_locks" table
CREATE TABLE `dahua_file_scan_locks` (`device_id` integer NOT NULL, `touched_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
CREATE UNIQUE INDEX `dahua_file_scan_locks_device_id` ON `dahua_file_scan_locks` (`device_id`);
