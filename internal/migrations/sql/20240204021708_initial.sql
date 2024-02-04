-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_events" table
CREATE TABLE `new_events` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `action` text NOT NULL, `slug` text NOT NULL, `actor` text NOT NULL, `user_id` integer NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE SET NULL);
-- copy rows from old table "events" to new temporary table "new_events"
INSERT INTO `new_events` (`id`, `action`, `slug`, `created_at`) SELECT `id`, `action`, `slug`, `created_at` FROM `events`;
-- drop "events" table after copying rows
DROP TABLE `events`;
-- rename temporary table "new_events" to "events"
ALTER TABLE `new_events` RENAME TO `events`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create "new_events" table
DROP TABLE `new_events`;
