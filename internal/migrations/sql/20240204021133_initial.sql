-- +goose Up
-- create "events" table
CREATE TABLE `events` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `action` text NOT NULL, `slug` text NOT NULL, `created_at` datetime NOT NULL);

-- +goose Down
-- reverse: create "events" table
DROP TABLE `events`;
