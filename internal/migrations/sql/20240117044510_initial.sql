-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_settings" table
CREATE TABLE `new_settings` (`setup` boolean NOT NULL, `site_name` text NOT NULL, `location` text NOT NULL, `coordinates` text NOT NULL);
-- copy rows from old table "settings" to new temporary table "new_settings"
INSERT INTO `new_settings` (`site_name`) SELECT `site_name` FROM `settings`;
-- drop "settings" table after copying rows
DROP TABLE `settings`;
-- rename temporary table "new_settings" to "settings"
ALTER TABLE `new_settings` RENAME TO `settings`;
-- create "users" table
CREATE TABLE `users` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `email` text NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL);
-- create index "users_email" to table: "users"
CREATE UNIQUE INDEX `users_email` ON `users` (`email`);
-- create index "users_username" to table: "users"
CREATE UNIQUE INDEX `users_username` ON `users` (`username`);
-- create "admins" table
CREATE TABLE `admins` (`user_id` integer NOT NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create "groups" table
CREATE TABLE `groups` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `description` text NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL);
-- create "group_users" table
CREATE TABLE `group_users` (`user_id` integer NOT NULL, `group_id` integer NOT NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create "dahua_permissions" table
CREATE TABLE `dahua_permissions` (`user_id` integer NULL, `group_id` integer NULL, `device_id` integer NOT NULL, `read` boolean NOT NULL, `write` boolean NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT `2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_permissions_user_id_device_id" to table: "dahua_permissions"
CREATE UNIQUE INDEX `dahua_permissions_user_id_device_id` ON `dahua_permissions` (`user_id`, `device_id`);
-- create index "dahua_permissions_group_id_device_id" to table: "dahua_permissions"
CREATE UNIQUE INDEX `dahua_permissions_group_id_device_id` ON `dahua_permissions` (`group_id`, `device_id`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "dahua_permissions_group_id_device_id" to table: "dahua_permissions"
DROP INDEX `dahua_permissions_group_id_device_id`;
-- reverse: create index "dahua_permissions_user_id_device_id" to table: "dahua_permissions"
DROP INDEX `dahua_permissions_user_id_device_id`;
-- reverse: create "dahua_permissions" table
DROP TABLE `dahua_permissions`;
-- reverse: create "group_users" table
DROP TABLE `group_users`;
-- reverse: create "groups" table
DROP TABLE `groups`;
-- reverse: create "admins" table
DROP TABLE `admins`;
-- reverse: create index "users_username" to table: "users"
DROP INDEX `users_username`;
-- reverse: create index "users_email" to table: "users"
DROP INDEX `users_email`;
-- reverse: create "users" table
DROP TABLE `users`;
-- reverse: create "new_settings" table
DROP TABLE `new_settings`;
