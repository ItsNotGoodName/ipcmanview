-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_dahua_files" table
CREATE TABLE `new_dahua_files` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `channel` integer NOT NULL, `start_time` datetime NOT NULL, `end_time` datetime NOT NULL, `length` integer NOT NULL, `type` text NOT NULL, `file_path` text NOT NULL, `duration` integer NOT NULL, `disk` integer NOT NULL, `video_stream` text NOT NULL, `flags` json NOT NULL, `events` json NOT NULL, `cluster` integer NOT NULL, `partition` integer NOT NULL, `pic_index` integer NOT NULL, `repeat` integer NOT NULL, `work_dir` text NOT NULL, `work_dir_sn` boolean NOT NULL, `updated_at` datetime NOT NULL, `storage` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- copy rows from old table "dahua_files" to new temporary table "new_dahua_files"
INSERT INTO `new_dahua_files` (`id`, `device_id`, `channel`, `start_time`, `end_time`, `length`, `type`, `file_path`, `duration`, `disk`, `video_stream`, `flags`, `events`, `cluster`, `partition`, `pic_index`, `repeat`, `work_dir`, `work_dir_sn`, `updated_at`, `storage`) SELECT `id`, `device_id`, `channel`, `start_time`, `end_time`, `length`, `type`, `file_path`, `duration`, `disk`, `video_stream`, `flags`, `events`, `cluster`, `partition`, `pic_index`, `repeat`, `work_dir`, `work_dir_sn`, `updated_at`, "local" FROM `dahua_files`;
-- drop "dahua_files" table after copying rows
DROP TABLE `dahua_files`;
-- rename temporary table "new_dahua_files" to "dahua_files"
ALTER TABLE `new_dahua_files` RENAME TO `dahua_files`;
-- create index "dahua_files_start_time" to table: "dahua_files"
CREATE UNIQUE INDEX `dahua_files_start_time` ON `dahua_files` (`start_time`);
-- create index "dahua_files_device_id_file_path" to table: "dahua_files"
CREATE UNIQUE INDEX `dahua_files_device_id_file_path` ON `dahua_files` (`device_id`, `file_path`);
-- create "new_dahua_credentials" table
CREATE TABLE `new_dahua_credentials` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `storage` text NOT NULL, `server_address` text NOT NULL, `port` integer NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `remote_directory` text NOT NULL);
-- copy rows from old table "dahua_credentials" to new temporary table "new_dahua_credentials"
INSERT INTO `new_dahua_credentials` (`id`, `server_address`, `port`, `username`, `password`, `remote_directory`) SELECT `id`, `server_address`, `port`, `username`, `password`, `remote_directory` FROM `dahua_credentials`;
-- drop "dahua_credentials" table after copying rows
DROP TABLE `dahua_credentials`;
-- rename temporary table "new_dahua_credentials" to "dahua_credentials"
ALTER TABLE `new_dahua_credentials` RENAME TO `dahua_credentials`;
-- create "new_dahua_device_credentials" table
CREATE TABLE `new_dahua_device_credentials` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `storage` text NOT NULL, `server_address` text NOT NULL, `port` integer NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `remote_directory` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- copy rows from old table "dahua_device_credentials" to new temporary table "new_dahua_device_credentials"
INSERT INTO `new_dahua_device_credentials` (`id`, `device_id`, `server_address`, `port`, `username`, `password`, `remote_directory`) SELECT `id`, `device_id`, `server_address`, `port`, `username`, `password`, `remote_directory` FROM `dahua_device_credentials`;
-- drop "dahua_device_credentials" table after copying rows
DROP TABLE `dahua_device_credentials`;
-- rename temporary table "new_dahua_device_credentials" to "dahua_device_credentials"
ALTER TABLE `new_dahua_device_credentials` RENAME TO `dahua_device_credentials`;
-- create index "dahua_device_credentials_storage" to table: "dahua_device_credentials"
CREATE UNIQUE INDEX `dahua_device_credentials_storage` ON `dahua_device_credentials` (`storage`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "dahua_device_credentials_storage" to table: "dahua_device_credentials"
DROP INDEX `dahua_device_credentials_storage`;
-- reverse: create "new_dahua_device_credentials" table
DROP TABLE `new_dahua_device_credentials`;
-- reverse: create "new_dahua_credentials" table
DROP TABLE `new_dahua_credentials`;
-- reverse: create index "dahua_files_device_id_file_path" to table: "dahua_files"
DROP INDEX `dahua_files_device_id_file_path`;
-- reverse: create index "dahua_files_start_time" to table: "dahua_files"
DROP INDEX `dahua_files_start_time`;
-- reverse: create "new_dahua_files" table
DROP TABLE `new_dahua_files`;
