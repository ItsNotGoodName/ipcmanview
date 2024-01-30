-- +goose Up
-- create "settings" table
CREATE TABLE `settings` (`setup` boolean NOT NULL, `site_name` text NOT NULL, `location` text NOT NULL, `coordinates` text NOT NULL);
-- create "users" table
CREATE TABLE `users` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `email` text NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `disabled_at` datetime NULL);
-- create index "users_email" to table: "users"
CREATE UNIQUE INDEX `users_email` ON `users` (`email`);
-- create index "users_username" to table: "users"
CREATE UNIQUE INDEX `users_username` ON `users` (`username`);
-- create "user_sessions" table
CREATE TABLE `user_sessions` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `user_id` integer NOT NULL, `session` text NOT NULL, `user_agent` text NOT NULL, `ip` text NOT NULL, `last_ip` text NOT NULL, `last_used_at` datetime NOT NULL, `created_at` datetime NOT NULL, `expired_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create "admins" table
CREATE TABLE `admins` (`user_id` integer NOT NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create "groups" table
CREATE TABLE `groups` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `description` text NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `disabled_at` datetime NULL);
-- create index "groups_name" to table: "groups"
CREATE UNIQUE INDEX `groups_name` ON `groups` (`name`);
-- create "group_users" table
CREATE TABLE `group_users` (`user_id` integer NOT NULL, `group_id` integer NOT NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "group_users_user_id_group_id" to table: "group_users"
CREATE UNIQUE INDEX `group_users_user_id_group_id` ON `group_users` (`user_id`, `group_id`);
-- create "dahua_devices" table
CREATE TABLE `dahua_devices` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `ip` text NOT NULL, `url` text NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `location` text NOT NULL, `feature` integer NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `disabled_at` datetime NULL);
-- create index "dahua_devices_name" to table: "dahua_devices"
CREATE UNIQUE INDEX `dahua_devices_name` ON `dahua_devices` (`name`);
-- create index "dahua_devices_ip" to table: "dahua_devices"
CREATE UNIQUE INDEX `dahua_devices_ip` ON `dahua_devices` (`ip`);
-- create "dahua_permissions" table
CREATE TABLE `dahua_permissions` (`user_id` integer NULL, `group_id` integer NULL, `device_id` integer NOT NULL, `level` integer NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT `2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_permissions_user_id_device_id" to table: "dahua_permissions"
CREATE UNIQUE INDEX `dahua_permissions_user_id_device_id` ON `dahua_permissions` (`user_id`, `device_id`);
-- create index "dahua_permissions_group_id_device_id" to table: "dahua_permissions"
CREATE UNIQUE INDEX `dahua_permissions_group_id_device_id` ON `dahua_permissions` (`group_id`, `device_id`);
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
-- create "dahua_event_worker_states" table
CREATE TABLE `dahua_event_worker_states` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `state` text NOT NULL, `error` text NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create "dahua_afero_files" table
CREATE TABLE `dahua_afero_files` (`id` integer NOT NULL, `file_id` integer NULL, `thumbnail_id` integer NULL, `email_attachment_id` integer NULL, `name` text NOT NULL, `ready` boolean NOT NULL DEFAULT false, `size` integer NOT NULL DEFAULT 0, `created_at` datetime NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`email_attachment_id`) REFERENCES `dahua_email_attachments` (`id`) ON UPDATE CASCADE ON DELETE SET NULL, CONSTRAINT `1` FOREIGN KEY (`thumbnail_id`) REFERENCES `dahua_thumbnails` (`id`) ON UPDATE CASCADE ON DELETE SET NULL, CONSTRAINT `2` FOREIGN KEY (`file_id`) REFERENCES `dahua_files` (`id`) ON UPDATE CASCADE ON DELETE SET NULL);
-- create index "dahua_afero_files_file_id" to table: "dahua_afero_files"
CREATE UNIQUE INDEX `dahua_afero_files_file_id` ON `dahua_afero_files` (`file_id`);
-- create index "dahua_afero_files_thumbnail_id" to table: "dahua_afero_files"
CREATE UNIQUE INDEX `dahua_afero_files_thumbnail_id` ON `dahua_afero_files` (`thumbnail_id`);
-- create index "dahua_afero_files_email_attachment_id" to table: "dahua_afero_files"
CREATE UNIQUE INDEX `dahua_afero_files_email_attachment_id` ON `dahua_afero_files` (`email_attachment_id`);
-- create index "dahua_afero_files_name" to table: "dahua_afero_files"
CREATE UNIQUE INDEX `dahua_afero_files_name` ON `dahua_afero_files` (`name`);
-- create "dahua_files" table
CREATE TABLE `dahua_files` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `channel` integer NOT NULL, `start_time` datetime NOT NULL, `end_time` datetime NOT NULL, `length` integer NOT NULL, `type` text NOT NULL, `file_path` text NOT NULL, `duration` integer NOT NULL, `disk` integer NOT NULL, `video_stream` text NOT NULL, `flags` json NOT NULL, `events` json NOT NULL, `cluster` integer NOT NULL, `partition` integer NOT NULL, `pic_index` integer NOT NULL, `repeat` integer NOT NULL, `work_dir` text NOT NULL, `work_dir_sn` boolean NOT NULL, `updated_at` datetime NOT NULL, `storage` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_files_start_time" to table: "dahua_files"
CREATE UNIQUE INDEX `dahua_files_start_time` ON `dahua_files` (`start_time`);
-- create index "dahua_files_device_id_file_path" to table: "dahua_files"
CREATE UNIQUE INDEX `dahua_files_device_id_file_path` ON `dahua_files` (`device_id`, `file_path`);
-- create "dahua_thumbnails" table
CREATE TABLE `dahua_thumbnails` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `file_id` integer NULL, `email_attachment_id` integer NULL, `width` integer NOT NULL, `height` integer NOT NULL, CONSTRAINT `0` FOREIGN KEY (`email_attachment_id`) REFERENCES `dahua_email_attachments` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`file_id`) REFERENCES `dahua_files` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_thumbnails_file_id_width_height" to table: "dahua_thumbnails"
CREATE UNIQUE INDEX `dahua_thumbnails_file_id_width_height` ON `dahua_thumbnails` (`file_id`, `width`, `height`);
-- create index "dahua_thumbnails_email_attachment_id_width_height" to table: "dahua_thumbnails"
CREATE UNIQUE INDEX `dahua_thumbnails_email_attachment_id_width_height` ON `dahua_thumbnails` (`email_attachment_id`, `width`, `height`);
-- create "dahua_file_cursors" table
CREATE TABLE `dahua_file_cursors` (`device_id` integer NOT NULL, `quick_cursor` datetime NOT NULL, `full_cursor` datetime NOT NULL, `full_epoch` datetime NOT NULL, `full_complete` boolean NOT NULL AS (full_cursor <= full_epoch) STORED, `scan` boolean NOT NULL, `scan_percent` real NOT NULL, `scan_type` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_file_cursors_device_id" to table: "dahua_file_cursors"
CREATE UNIQUE INDEX `dahua_file_cursors_device_id` ON `dahua_file_cursors` (`device_id`);
-- create "dahua_storage_destinations" table
CREATE TABLE `dahua_storage_destinations` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `storage` text NOT NULL, `server_address` text NOT NULL, `port` integer NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `remote_directory` text NOT NULL);
-- create index "dahua_storage_destinations_name" to table: "dahua_storage_destinations"
CREATE UNIQUE INDEX `dahua_storage_destinations_name` ON `dahua_storage_destinations` (`name`);
-- create "dahua_streams" table
CREATE TABLE `dahua_streams` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `internal` boolean NOT NULL, `device_id` integer NOT NULL, `channel` integer NOT NULL, `subtype` integer NOT NULL, `name` text NOT NULL, `mediamtx_path` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_streams_device_id_channel_subtype" to table: "dahua_streams"
CREATE UNIQUE INDEX `dahua_streams_device_id_channel_subtype` ON `dahua_streams` (`device_id`, `channel`, `subtype`);
-- create "dahua_email_messages" table
CREATE TABLE `dahua_email_messages` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `date` datetime NOT NULL, `from` text NOT NULL, `to` json NOT NULL, `subject` text NOT NULL, `text` text NOT NULL, `alarm_event` text NOT NULL, `alarm_input_channel` integer NOT NULL, `alarm_name` text NOT NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create "dahua_email_attachments" table
CREATE TABLE `dahua_email_attachments` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `message_id` integer NOT NULL, `file_name` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`message_id`) REFERENCES `dahua_email_messages` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);

-- +goose Down
-- reverse: create "dahua_email_attachments" table
DROP TABLE `dahua_email_attachments`;
-- reverse: create "dahua_email_messages" table
DROP TABLE `dahua_email_messages`;
-- reverse: create index "dahua_streams_device_id_channel_subtype" to table: "dahua_streams"
DROP INDEX `dahua_streams_device_id_channel_subtype`;
-- reverse: create "dahua_streams" table
DROP TABLE `dahua_streams`;
-- reverse: create index "dahua_storage_destinations_name" to table: "dahua_storage_destinations"
DROP INDEX `dahua_storage_destinations_name`;
-- reverse: create "dahua_storage_destinations" table
DROP TABLE `dahua_storage_destinations`;
-- reverse: create index "dahua_file_cursors_device_id" to table: "dahua_file_cursors"
DROP INDEX `dahua_file_cursors_device_id`;
-- reverse: create "dahua_file_cursors" table
DROP TABLE `dahua_file_cursors`;
-- reverse: create index "dahua_thumbnails_email_attachment_id_width_height" to table: "dahua_thumbnails"
DROP INDEX `dahua_thumbnails_email_attachment_id_width_height`;
-- reverse: create index "dahua_thumbnails_file_id_width_height" to table: "dahua_thumbnails"
DROP INDEX `dahua_thumbnails_file_id_width_height`;
-- reverse: create "dahua_thumbnails" table
DROP TABLE `dahua_thumbnails`;
-- reverse: create index "dahua_files_device_id_file_path" to table: "dahua_files"
DROP INDEX `dahua_files_device_id_file_path`;
-- reverse: create index "dahua_files_start_time" to table: "dahua_files"
DROP INDEX `dahua_files_start_time`;
-- reverse: create "dahua_files" table
DROP TABLE `dahua_files`;
-- reverse: create index "dahua_afero_files_name" to table: "dahua_afero_files"
DROP INDEX `dahua_afero_files_name`;
-- reverse: create index "dahua_afero_files_email_attachment_id" to table: "dahua_afero_files"
DROP INDEX `dahua_afero_files_email_attachment_id`;
-- reverse: create index "dahua_afero_files_thumbnail_id" to table: "dahua_afero_files"
DROP INDEX `dahua_afero_files_thumbnail_id`;
-- reverse: create index "dahua_afero_files_file_id" to table: "dahua_afero_files"
DROP INDEX `dahua_afero_files_file_id`;
-- reverse: create "dahua_afero_files" table
DROP TABLE `dahua_afero_files`;
-- reverse: create "dahua_event_worker_states" table
DROP TABLE `dahua_event_worker_states`;
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
-- reverse: create index "dahua_permissions_group_id_device_id" to table: "dahua_permissions"
DROP INDEX `dahua_permissions_group_id_device_id`;
-- reverse: create index "dahua_permissions_user_id_device_id" to table: "dahua_permissions"
DROP INDEX `dahua_permissions_user_id_device_id`;
-- reverse: create "dahua_permissions" table
DROP TABLE `dahua_permissions`;
-- reverse: create index "dahua_devices_ip" to table: "dahua_devices"
DROP INDEX `dahua_devices_ip`;
-- reverse: create index "dahua_devices_name" to table: "dahua_devices"
DROP INDEX `dahua_devices_name`;
-- reverse: create "dahua_devices" table
DROP TABLE `dahua_devices`;
-- reverse: create index "group_users_user_id_group_id" to table: "group_users"
DROP INDEX `group_users_user_id_group_id`;
-- reverse: create "group_users" table
DROP TABLE `group_users`;
-- reverse: create index "groups_name" to table: "groups"
DROP INDEX `groups_name`;
-- reverse: create "groups" table
DROP TABLE `groups`;
-- reverse: create "admins" table
DROP TABLE `admins`;
-- reverse: create "user_sessions" table
DROP TABLE `user_sessions`;
-- reverse: create index "users_username" to table: "users"
DROP INDEX `users_username`;
-- reverse: create index "users_email" to table: "users"
DROP INDEX `users_email`;
-- reverse: create "users" table
DROP TABLE `users`;
-- reverse: create "settings" table
DROP TABLE `settings`;
