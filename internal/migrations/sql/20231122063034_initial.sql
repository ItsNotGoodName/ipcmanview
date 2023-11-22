-- +goose Up
-- create "settings" table
CREATE TABLE `settings` (`site_name` text NOT NULL, `default_location` text NOT NULL);
-- create "dahua_cameras" table
CREATE TABLE `dahua_cameras` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `address` text NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `location` text NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL);
-- create index "dahua_cameras_name" to table: "dahua_cameras"
CREATE UNIQUE INDEX `dahua_cameras_name` ON `dahua_cameras` (`name`);
-- create index "dahua_cameras_address" to table: "dahua_cameras"
CREATE UNIQUE INDEX `dahua_cameras_address` ON `dahua_cameras` (`address`);
-- create "dahua_seeds" table
CREATE TABLE `dahua_seeds` (`seed` integer NOT NULL, `camera_id` integer NULL, PRIMARY KEY (`seed`), CONSTRAINT `0` FOREIGN KEY (`camera_id`) REFERENCES `dahua_cameras` (`id`) ON UPDATE CASCADE ON DELETE SET NULL);
-- create index "dahua_seeds_camera_id" to table: "dahua_seeds"
CREATE UNIQUE INDEX `dahua_seeds_camera_id` ON `dahua_seeds` (`camera_id`);
-- create "dahua_events" table
CREATE TABLE `dahua_events` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `camera_id` integer NOT NULL, `content_type` text NOT NULL, `content_length` integer NOT NULL, `code` text NOT NULL, `action` text NOT NULL, `index` integer NOT NULL, `data` json NOT NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`camera_id`) REFERENCES `dahua_cameras` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create "dahua_files" table
CREATE TABLE `dahua_files` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `camera_id` integer NOT NULL, `file_path` text NOT NULL, `kind` text NOT NULL, `size` integer NOT NULL, `start_time` datetime NOT NULL, `end_time` datetime NOT NULL, `duration` integer NOT NULL, `events` json NOT NULL, `updated_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`camera_id`) REFERENCES `dahua_cameras` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_files_start_time" to table: "dahua_files"
CREATE UNIQUE INDEX `dahua_files_start_time` ON `dahua_files` (`start_time`);
-- create index "dahua_files_camera_id_file_path" to table: "dahua_files"
CREATE UNIQUE INDEX `dahua_files_camera_id_file_path` ON `dahua_files` (`camera_id`, `file_path`);
-- create "dahua_file_cursors" table
CREATE TABLE `dahua_file_cursors` (`camera_id` integer NOT NULL, `quick_cursor` datetime NOT NULL, `full_cursor` datetime NOT NULL, `full_epoch` datetime NOT NULL, `full_epoch_end` datetime NOT NULL, `full_complete` boolean NOT NULL AS (full_cursor <= full_epoch) STORED, CONSTRAINT `0` FOREIGN KEY (`camera_id`) REFERENCES `dahua_cameras` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CHECK (full_cursor <= full_epoch_end));
-- create index "dahua_file_cursors_camera_id" to table: "dahua_file_cursors"
CREATE UNIQUE INDEX `dahua_file_cursors_camera_id` ON `dahua_file_cursors` (`camera_id`);
-- create "dahua_file_scan_locks" table
CREATE TABLE `dahua_file_scan_locks` (`camera_id` integer NOT NULL, `created_at` datetime NOT NULL);
-- create index "dahua_file_scan_locks_camera_id" to table: "dahua_file_scan_locks"
CREATE UNIQUE INDEX `dahua_file_scan_locks_camera_id` ON `dahua_file_scan_locks` (`camera_id`);

WITH RECURSIVE generate_series(value) AS (
  SELECT 1
  UNION ALL
  SELECT value+1 FROM generate_series WHERE value+1<=999
)
INSERT INTO dahua_seeds (seed) SELECT value from generate_series;

-- +goose Down
-- reverse: create index "dahua_file_scan_locks_camera_id" to table: "dahua_file_scan_locks"
DROP INDEX `dahua_file_scan_locks_camera_id`;
-- reverse: create "dahua_file_scan_locks" table
DROP TABLE `dahua_file_scan_locks`;
-- reverse: create index "dahua_file_cursors_camera_id" to table: "dahua_file_cursors"
DROP INDEX `dahua_file_cursors_camera_id`;
-- reverse: create "dahua_file_cursors" table
DROP TABLE `dahua_file_cursors`;
-- reverse: create index "dahua_files_camera_id_file_path" to table: "dahua_files"
DROP INDEX `dahua_files_camera_id_file_path`;
-- reverse: create index "dahua_files_start_time" to table: "dahua_files"
DROP INDEX `dahua_files_start_time`;
-- reverse: create "dahua_files" table
DROP TABLE `dahua_files`;
-- reverse: create "dahua_events" table
DROP TABLE `dahua_events`;
-- reverse: create index "dahua_seeds_camera_id" to table: "dahua_seeds"
DROP INDEX `dahua_seeds_camera_id`;
-- reverse: create "dahua_seeds" table
DROP TABLE `dahua_seeds`;
-- reverse: create index "dahua_cameras_address" to table: "dahua_cameras"
DROP INDEX `dahua_cameras_address`;
-- reverse: create index "dahua_cameras_name" to table: "dahua_cameras"
DROP INDEX `dahua_cameras_name`;
-- reverse: create "dahua_cameras" table
DROP TABLE `dahua_cameras`;
-- reverse: create "settings" table
DROP TABLE `settings`;

WITH RECURSIVE generate_series(value) AS (
  SELECT 1
  UNION ALL
  SELECT value+1 FROM generate_series WHERE value+1<=999
)
DELETE FROM dahua_seeds WHERE seed IN (SELECT value from generate_series);
