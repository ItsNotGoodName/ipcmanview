-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_dahua_streams" table
CREATE TABLE `new_dahua_streams` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `channel` integer NOT NULL, `subtype` integer NOT NULL, `name` text NOT NULL, `mediamtx_path` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- copy rows from old table "dahua_streams" to new temporary table "new_dahua_streams"
INSERT INTO `new_dahua_streams` (`id`, `device_id`, `channel`, `subtype`, `name`, `mediamtx_path`) SELECT `id`, `device_id`, `channel`, `subtype`, `name`, `mediamtx_path` FROM `dahua_streams`;
-- drop "dahua_streams" table after copying rows
DROP TABLE `dahua_streams`;
-- rename temporary table "new_dahua_streams" to "dahua_streams"
ALTER TABLE `new_dahua_streams` RENAME TO `dahua_streams`;
-- create index "dahua_streams_device_id_channel_subtype" to table: "dahua_streams"
CREATE UNIQUE INDEX `dahua_streams_device_id_channel_subtype` ON `dahua_streams` (`device_id`, `channel`, `subtype`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "dahua_streams_device_id_channel_subtype" to table: "dahua_streams"
DROP INDEX `dahua_streams_device_id_channel_subtype`;
-- reverse: create "new_dahua_streams" table
DROP TABLE `new_dahua_streams`;
