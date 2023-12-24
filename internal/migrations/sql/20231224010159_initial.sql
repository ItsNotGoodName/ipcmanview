-- +goose Up
-- create "settings" table
CREATE TABLE `settings` (`site_name` text NOT NULL, `default_location` text NOT NULL);
-- create "dahua_devices" table
CREATE TABLE `dahua_devices` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `address` text NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `location` text NOT NULL, `feature` integer NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL);
-- create index "dahua_devices_name" to table: "dahua_devices"
CREATE UNIQUE INDEX `dahua_devices_name` ON `dahua_devices` (`name`);
-- create index "dahua_devices_address" to table: "dahua_devices"
CREATE UNIQUE INDEX `dahua_devices_address` ON `dahua_devices` (`address`);
-- create "dahua_seeds" table
CREATE TABLE `dahua_seeds` (`seed` integer NOT NULL, `device_id` integer NULL, PRIMARY KEY (`seed`), CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE SET NULL);
-- create index "dahua_seeds_device_id" to table: "dahua_seeds"
CREATE UNIQUE INDEX `dahua_seeds_device_id` ON `dahua_seeds` (`device_id`);
-- create "dahua_events" table
CREATE TABLE `dahua_events` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `code` text NOT NULL, `action` text NOT NULL, `index` integer NOT NULL, `data` json NOT NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create "dahua_event_rules" table
CREATE TABLE `dahua_event_rules` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `code` text NOT NULL, `ignore_db` boolean NOT NULL DEFAULT false, `ignore_live` boolean NOT NULL DEFAULT false, `ignore_mqtt` boolean NOT NULL DEFAULT false);
-- create index "dahua_event_rules_code" to table: "dahua_event_rules"
CREATE UNIQUE INDEX `dahua_event_rules_code` ON `dahua_event_rules` (`code`);
-- create "dahua_event_device_rules" table
CREATE TABLE `dahua_event_device_rules` (`device_id` integer NOT NULL, `code` text NOT NULL, `ignore_db` boolean NOT NULL DEFAULT false, `ignore_live` boolean NOT NULL DEFAULT false, `ignore_mqtt` boolean NOT NULL DEFAULT false, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_event_device_rules_device_id_code" to table: "dahua_event_device_rules"
CREATE UNIQUE INDEX `dahua_event_device_rules_device_id_code` ON `dahua_event_device_rules` (`device_id`, `code`);
-- create "dahua_files" table
CREATE TABLE `dahua_files` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `channel` integer NOT NULL, `start_time` datetime NOT NULL, `end_time` datetime NOT NULL, `length` integer NOT NULL, `type` text NOT NULL, `file_path` text NOT NULL, `duration` integer NOT NULL, `disk` integer NOT NULL, `video_stream` text NOT NULL, `flags` json NOT NULL, `events` json NOT NULL, `cluster` integer NOT NULL, `partition` integer NOT NULL, `pic_index` integer NOT NULL, `repeat` integer NOT NULL, `work_dir` text NOT NULL, `work_dir_sn` boolean NOT NULL, `updated_at` datetime NOT NULL, `local` boolean NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_files_start_time" to table: "dahua_files"
CREATE UNIQUE INDEX `dahua_files_start_time` ON `dahua_files` (`start_time`);
-- create index "dahua_files_device_id_file_path" to table: "dahua_files"
CREATE UNIQUE INDEX `dahua_files_device_id_file_path` ON `dahua_files` (`device_id`, `file_path`);
-- create "dahua_file_cursors" table
CREATE TABLE `dahua_file_cursors` (`device_id` integer NOT NULL, `quick_cursor` datetime NOT NULL, `full_cursor` datetime NOT NULL, `full_epoch` datetime NOT NULL, `full_complete` boolean NOT NULL AS (full_cursor <= full_epoch) STORED, `percent` real NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_file_cursors_device_id" to table: "dahua_file_cursors"
CREATE UNIQUE INDEX `dahua_file_cursors_device_id` ON `dahua_file_cursors` (`device_id`);
-- create "dahua_file_scan_locks" table
CREATE TABLE `dahua_file_scan_locks` (`device_id` integer NOT NULL, `touched_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_file_scan_locks_device_id" to table: "dahua_file_scan_locks"
CREATE UNIQUE INDEX `dahua_file_scan_locks_device_id` ON `dahua_file_scan_locks` (`device_id`);
-- create "dahua_event_worker_states" table
CREATE TABLE `dahua_event_worker_states` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `state` text NOT NULL, `error` text NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);

WITH RECURSIVE generate_series(value) AS (
  SELECT 1
  UNION ALL
  SELECT value+1 FROM generate_series WHERE value+1<=999
)
INSERT INTO dahua_seeds (seed) SELECT value from generate_series;
INSERT INTO dahua_event_rules (code) VALUES ('');

-- +goose Down
-- reverse: create "dahua_event_worker_states" table
DROP TABLE `dahua_event_worker_states`;
-- reverse: create index "dahua_file_scan_locks_device_id" to table: "dahua_file_scan_locks"
DROP INDEX `dahua_file_scan_locks_device_id`;
-- reverse: create "dahua_file_scan_locks" table
DROP TABLE `dahua_file_scan_locks`;
-- reverse: create index "dahua_file_cursors_device_id" to table: "dahua_file_cursors"
DROP INDEX `dahua_file_cursors_device_id`;
-- reverse: create "dahua_file_cursors" table
DROP TABLE `dahua_file_cursors`;
-- reverse: create index "dahua_files_device_id_file_path" to table: "dahua_files"
DROP INDEX `dahua_files_device_id_file_path`;
-- reverse: create index "dahua_files_start_time" to table: "dahua_files"
DROP INDEX `dahua_files_start_time`;
-- reverse: create "dahua_files" table
DROP TABLE `dahua_files`;
-- reverse: create index "dahua_event_device_rules_device_id_code" to table: "dahua_event_device_rules"
DROP INDEX `dahua_event_device_rules_device_id_code`;
-- reverse: create "dahua_event_device_rules" table
DROP TABLE `dahua_event_device_rules`;
-- reverse: create index "dahua_event_rules_code" to table: "dahua_event_rules"
DROP INDEX `dahua_event_rules_code`;
-- reverse: create "dahua_event_rules" table
DROP TABLE `dahua_event_rules`;
-- reverse: create "dahua_events" table
DROP TABLE `dahua_events`;
-- reverse: create index "dahua_seeds_device_id" to table: "dahua_seeds"
DROP INDEX `dahua_seeds_device_id`;
-- reverse: create "dahua_seeds" table
DROP TABLE `dahua_seeds`;
-- reverse: create index "dahua_devices_address" to table: "dahua_devices"
DROP INDEX `dahua_devices_address`;
-- reverse: create index "dahua_devices_name" to table: "dahua_devices"
DROP INDEX `dahua_devices_name`;
-- reverse: create "dahua_devices" table
DROP TABLE `dahua_devices`;
-- reverse: create "settings" table
DROP TABLE `settings`;

WITH RECURSIVE generate_series(value) AS (
  SELECT 1
  UNION ALL
  SELECT value+1 FROM generate_series WHERE value+1<=999
)
DELETE FROM dahua_seeds WHERE seed IN (SELECT value from generate_series);
DELETE FROM dahua_event_rules WHERE code = '';
