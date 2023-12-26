-- +goose Up
-- create "dahua_credentials" table
CREATE TABLE `dahua_credentials` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `type` text NOT NULL, `server_address` text NOT NULL, `port` integer NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `remote_directory` text NOT NULL);
-- create "dahua_device_credentials" table
CREATE TABLE `dahua_device_credentials` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `type` text NOT NULL, `server_address` text NOT NULL, `port` integer NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `remote_directory` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_device_credentials_type" to table: "dahua_device_credentials"
CREATE UNIQUE INDEX `dahua_device_credentials_type` ON `dahua_device_credentials` (`type`);

-- +goose Down
-- reverse: create index "dahua_device_credentials_type" to table: "dahua_device_credentials"
DROP INDEX `dahua_device_credentials_type`;
-- reverse: create "dahua_device_credentials" table
DROP TABLE `dahua_device_credentials`;
-- reverse: create "dahua_credentials" table
DROP TABLE `dahua_credentials`;
