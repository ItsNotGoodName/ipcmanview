-- +goose Up
-- create "user_sessions" table
CREATE TABLE `user_sessions` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `user_id` integer NOT NULL, `session` text NOT NULL, `created_at` datetime NOT NULL, `expires_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);

-- +goose Down
-- reverse: create "user_sessions" table
DROP TABLE `user_sessions`;
