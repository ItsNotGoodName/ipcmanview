-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_dahua_permissions" table
CREATE TABLE `new_dahua_permissions` (`user_id` integer NULL, `group_id` integer NULL, `device_id` integer NOT NULL, `level` integer NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT `2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- copy rows from old table "dahua_permissions" to new temporary table "new_dahua_permissions"
INSERT INTO `new_dahua_permissions` (`user_id`, `group_id`, `device_id`) SELECT `user_id`, `group_id`, `device_id` FROM `dahua_permissions`;
-- drop "dahua_permissions" table after copying rows
DROP TABLE `dahua_permissions`;
-- rename temporary table "new_dahua_permissions" to "dahua_permissions"
ALTER TABLE `new_dahua_permissions` RENAME TO `dahua_permissions`;
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
-- reverse: create "new_dahua_permissions" table
DROP TABLE `new_dahua_permissions`;
