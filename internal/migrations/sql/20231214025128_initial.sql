-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_dahua_event_worker_states" table
CREATE TABLE `new_dahua_event_worker_states` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `camera_id` integer NOT NULL, `state` text NOT NULL, `error` text NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`camera_id`) REFERENCES `dahua_cameras` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- copy rows from old table "dahua_event_worker_states" to new temporary table "new_dahua_event_worker_states"
INSERT INTO `new_dahua_event_worker_states` (`id`, `camera_id`, `state`, `error`, `created_at`) SELECT `id`, `camera_id`, `action`, `error`, `created_at` FROM `dahua_event_worker_states`;
-- drop "dahua_event_worker_states" table after copying rows
DROP TABLE `dahua_event_worker_states`;
-- rename temporary table "new_dahua_event_worker_states" to "dahua_event_worker_states"
ALTER TABLE `new_dahua_event_worker_states` RENAME TO `dahua_event_worker_states`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create "new_dahua_event_worker_states" table
DROP TABLE `new_dahua_event_worker_states`;
