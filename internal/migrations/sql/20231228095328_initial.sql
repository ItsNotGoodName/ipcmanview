-- +goose Up
-- create "dahua_streams" table
CREATE TABLE `dahua_streams` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `channel` integer NOT NULL, `subtype` integer NOT NULL, `mediamtx_path` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_streams_device_id_channel_subtype" to table: "dahua_streams"
CREATE UNIQUE INDEX `dahua_streams_device_id_channel_subtype` ON `dahua_streams` (`device_id`, `channel`, `subtype`);

-- +goose Down
-- reverse: create index "dahua_streams_device_id_channel_subtype" to table: "dahua_streams"
DROP INDEX `dahua_streams_device_id_channel_subtype`;
-- reverse: create "dahua_streams" table
DROP TABLE `dahua_streams`;
