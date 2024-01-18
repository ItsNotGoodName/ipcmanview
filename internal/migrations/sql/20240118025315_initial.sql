-- +goose Up
-- create "user_tokens" table
CREATE TABLE `user_tokens` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `user_id` integer NOT NULL, `token` text NOT NULL, `created_at` datetime NOT NULL, `revoked_at` datetime NULL, CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);

-- +goose Down
-- reverse: create "user_tokens" table
DROP TABLE `user_tokens`;
