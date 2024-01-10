-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- drop "dahua_credentials" table
DROP TABLE `dahua_credentials`;
-- drop "dahua_device_credentials" table
DROP TABLE `dahua_device_credentials`;
-- create "dahua_storage_destinations" table
CREATE TABLE `dahua_storage_destinations` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `storage` text NOT NULL, `server_address` text NOT NULL, `port` integer NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `remote_directory` text NOT NULL);
-- create index "dahua_storage_destinations_name" to table: "dahua_storage_destinations"
CREATE UNIQUE INDEX `dahua_storage_destinations_name` ON `dahua_storage_destinations` (`name`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "dahua_storage_destinations_name" to table: "dahua_storage_destinations"
DROP INDEX `dahua_storage_destinations_name`;
-- reverse: create "dahua_storage_destinations" table
DROP TABLE `dahua_storage_destinations`;
-- reverse: drop "dahua_device_credentials" table
CREATE TABLE `dahua_device_credentials` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `storage` text NOT NULL, `server_address` text NOT NULL, `port` integer NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `remote_directory` text NOT NULL, `name` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
CREATE UNIQUE INDEX `dahua_device_credentials_storage` ON `dahua_device_credentials` (`storage`);
CREATE UNIQUE INDEX `dahua_device_credentials_name` ON `dahua_device_credentials` (`name`);
-- reverse: drop "dahua_credentials" table
CREATE TABLE `dahua_credentials` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `storage` text NOT NULL, `server_address` text NOT NULL, `port` integer NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `remote_directory` text NOT NULL, `name` text NOT NULL);
CREATE UNIQUE INDEX `dahua_credentials_name` ON `dahua_credentials` (`name`);
