-- +goose Up
-- create "dahua_cameras" table
CREATE TABLE `dahua_cameras` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `address` text NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `location` text NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL);
-- create index "dahua_cameras_name" to table: "dahua_cameras"
CREATE UNIQUE INDEX `dahua_cameras_name` ON `dahua_cameras` (`name`);
-- create index "dahua_cameras_address" to table: "dahua_cameras"
CREATE UNIQUE INDEX `dahua_cameras_address` ON `dahua_cameras` (`address`);
-- create "dahua_seeds" table
CREATE TABLE `dahua_seeds` (`seed` integer NOT NULL, `camera_id` integer NULL, PRIMARY KEY (`seed`), CONSTRAINT `0` FOREIGN KEY (`camera_id`) REFERENCES `dahua_cameras` (`id`) ON UPDATE CASCADE ON DELETE SET NULL);
-- create "dahua_events" table
CREATE TABLE `dahua_events` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `camera_id` integer NOT NULL, `content_type` text NOT NULL, `content_length` integer NOT NULL, `code` text NOT NULL, `action` text NOT NULL, `index` integer NOT NULL, `data` json NOT NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`camera_id`) REFERENCES `dahua_cameras` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create "dahua_camera_files" table
CREATE TABLE `dahua_camera_files` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `camera_id` integer NOT NULL, `file_path` text NOT NULL, `kind` text NOT NULL, `size` integer NOT NULL, `start_time` datetime NOT NULL, `end_time` datetime NOT NULL, `duration` integer NOT NULL, `events` json NOT NULL, `updated_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`camera_id`) REFERENCES `dahua_cameras` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_camera_files_start_time" to table: "dahua_camera_files"
CREATE UNIQUE INDEX `dahua_camera_files_start_time` ON `dahua_camera_files` (`start_time`);
-- create index "dahua_camera_files_camera_id_file_path" to table: "dahua_camera_files"
CREATE UNIQUE INDEX `dahua_camera_files_camera_id_file_path` ON `dahua_camera_files` (`camera_id`, `file_path`);

WITH RECURSIVE generate_series(value) AS (
  SELECT 1
  UNION ALL
  SELECT value+1 FROM generate_series WHERE value+1<=999
)
INSERT INTO dahua_seeds (seed) SELECT value from generate_series;

-- +goose Down
-- reverse: create index "dahua_camera_files_camera_id_file_path" to table: "dahua_camera_files"
DROP INDEX `dahua_camera_files_camera_id_file_path`;
-- reverse: create index "dahua_camera_files_start_time" to table: "dahua_camera_files"
DROP INDEX `dahua_camera_files_start_time`;
-- reverse: create "dahua_camera_files" table
DROP TABLE `dahua_camera_files`;
-- reverse: create "dahua_events" table
DROP TABLE `dahua_events`;
-- reverse: create "dahua_seeds" table
DROP TABLE `dahua_seeds`;
-- reverse: create index "dahua_cameras_address" to table: "dahua_cameras"
DROP INDEX `dahua_cameras_address`;
-- reverse: create index "dahua_cameras_name" to table: "dahua_cameras"
DROP INDEX `dahua_cameras_name`;
-- reverse: create "dahua_cameras" table
DROP TABLE `dahua_cameras`;

WITH RECURSIVE generate_series(value) AS (
  SELECT 1
  UNION ALL
  SELECT value+1 FROM generate_series WHERE value+1<=999
)
DELETE FROM dahua_seeds WHERE seed IN (SELECT value from generate_series);
