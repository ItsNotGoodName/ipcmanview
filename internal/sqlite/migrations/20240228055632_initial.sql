-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- drop "dahua_event_worker_states" table
DROP TABLE `dahua_event_worker_states`;
-- create "dahua_worker_events" table
CREATE TABLE `dahua_worker_events` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `type` text NOT NULL, `state` text NOT NULL, `error` text NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create "dahua_worker_events" table
DROP TABLE `dahua_worker_events`;
-- reverse: drop "dahua_event_worker_states" table
CREATE TABLE `dahua_event_worker_states` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `device_id` integer NOT NULL, `state` text NOT NULL, `error` text NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
