-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- add column "user_agent" to table: "user_sessions"
ALTER TABLE `user_sessions` ADD COLUMN `user_agent` text NOT NULL;
-- add column "ip" to table: "user_sessions"
ALTER TABLE `user_sessions` ADD COLUMN `ip` text NOT NULL;
-- add column "last_ip" to table: "user_sessions"
ALTER TABLE `user_sessions` ADD COLUMN `last_ip` text NOT NULL;
-- add column "last_used_at" to table: "user_sessions"
ALTER TABLE `user_sessions` ADD COLUMN `last_used_at` datetime NOT NULL;
-- drop "user_tokens" table
DROP TABLE `user_tokens`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: drop "user_tokens" table
CREATE TABLE `user_tokens` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `user_id` integer NOT NULL, `token` text NOT NULL, `created_at` datetime NOT NULL, `revoked_at` datetime NULL, CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- reverse: add column "last_used_at" to table: "user_sessions"
ALTER TABLE `user_sessions` DROP COLUMN `last_used_at`;
-- reverse: add column "last_ip" to table: "user_sessions"
ALTER TABLE `user_sessions` DROP COLUMN `last_ip`;
-- reverse: add column "ip" to table: "user_sessions"
ALTER TABLE `user_sessions` DROP COLUMN `ip`;
-- reverse: add column "user_agent" to table: "user_sessions"
ALTER TABLE `user_sessions` DROP COLUMN `user_agent`;
