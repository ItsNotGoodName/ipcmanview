-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_user_sessions" table
CREATE TABLE `new_user_sessions` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `user_id` integer NOT NULL, `session` text NOT NULL, `created_at` datetime NOT NULL, `expired_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- copy rows from old table "user_sessions" to new temporary table "new_user_sessions"
INSERT INTO `new_user_sessions` (`id`, `user_id`, `session`, `created_at`) SELECT `id`, `user_id`, `session`, `created_at` FROM `user_sessions`;
-- drop "user_sessions" table after copying rows
DROP TABLE `user_sessions`;
-- rename temporary table "new_user_sessions" to "user_sessions"
ALTER TABLE `new_user_sessions` RENAME TO `user_sessions`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create "new_user_sessions" table
DROP TABLE `new_user_sessions`;
