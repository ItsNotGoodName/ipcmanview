-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_dahua_event_default_rules" table
CREATE TABLE `new_dahua_event_default_rules` (`code` text NOT NULL DEFAULT '', `ignore_db` boolean NOT NULL DEFAULT false, `ignore_live` boolean NOT NULL DEFAULT false, `ignore_mqtt` boolean NOT NULL DEFAULT false);
-- copy rows from old table "dahua_event_default_rules" to new temporary table "new_dahua_event_default_rules"
INSERT INTO `new_dahua_event_default_rules` (`code`, `ignore_db`, `ignore_live`, `ignore_mqtt`) SELECT IFNULL(`code`, '') AS `code`, `ignore_db`, `ignore_live`, `ignore_mqtt` FROM `dahua_event_default_rules`;
-- drop "dahua_event_default_rules" table after copying rows
DROP TABLE `dahua_event_default_rules`;
-- rename temporary table "new_dahua_event_default_rules" to "dahua_event_default_rules"
ALTER TABLE `new_dahua_event_default_rules` RENAME TO `dahua_event_default_rules`;
-- create index "dahua_event_default_rules_code" to table: "dahua_event_default_rules"
CREATE UNIQUE INDEX `dahua_event_default_rules_code` ON `dahua_event_default_rules` (`code`);
-- create "new_dahua_event_rules" table
CREATE TABLE `new_dahua_event_rules` (`camera_id` integer NOT NULL, `code` text NOT NULL DEFAULT '', `ignore_db` boolean NOT NULL DEFAULT false, `ignore_live` boolean NOT NULL DEFAULT false, `ignore_mqtt` boolean NOT NULL DEFAULT false, CONSTRAINT `0` FOREIGN KEY (`camera_id`) REFERENCES `dahua_cameras` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- copy rows from old table "dahua_event_rules" to new temporary table "new_dahua_event_rules"
INSERT INTO `new_dahua_event_rules` (`camera_id`, `code`, `ignore_db`, `ignore_live`, `ignore_mqtt`) SELECT `camera_id`, IFNULL(`code`, '') AS `code`, `ignore_db`, `ignore_live`, `ignore_mqtt` FROM `dahua_event_rules`;
-- drop "dahua_event_rules" table after copying rows
DROP TABLE `dahua_event_rules`;
-- rename temporary table "new_dahua_event_rules" to "dahua_event_rules"
ALTER TABLE `new_dahua_event_rules` RENAME TO `dahua_event_rules`;
-- create index "dahua_event_rules_camera_id_code" to table: "dahua_event_rules"
CREATE UNIQUE INDEX `dahua_event_rules_camera_id_code` ON `dahua_event_rules` (`camera_id`, `code`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "dahua_event_rules_camera_id_code" to table: "dahua_event_rules"
DROP INDEX `dahua_event_rules_camera_id_code`;
-- reverse: create "new_dahua_event_rules" table
DROP TABLE `new_dahua_event_rules`;
-- reverse: create index "dahua_event_default_rules_code" to table: "dahua_event_default_rules"
DROP INDEX `dahua_event_default_rules_code`;
-- reverse: create "new_dahua_event_default_rules" table
DROP TABLE `new_dahua_event_default_rules`;
