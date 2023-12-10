-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_dahua_events" table
CREATE TABLE `new_dahua_events` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `camera_id` integer NOT NULL, `code` text NOT NULL, `action` text NOT NULL, `index` integer NOT NULL, `data` json NOT NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`camera_id`) REFERENCES `dahua_cameras` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- copy rows from old table "dahua_events" to new temporary table "new_dahua_events"
INSERT INTO `new_dahua_events` (`id`, `camera_id`, `code`, `action`, `index`, `data`, `created_at`) SELECT `id`, `camera_id`, `code`, `action`, `index`, `data`, `created_at` FROM `dahua_events`;
-- drop "dahua_events" table after copying rows
DROP TABLE `dahua_events`;
-- rename temporary table "new_dahua_events" to "dahua_events"
ALTER TABLE `new_dahua_events` RENAME TO `dahua_events`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create "new_dahua_events" table
DROP TABLE `new_dahua_events`;
