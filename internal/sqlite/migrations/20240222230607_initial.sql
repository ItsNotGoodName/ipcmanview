-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- drop "settings" table
DROP TABLE `settings`;
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: drop "settings" table
CREATE TABLE `settings` (`setup` boolean NOT NULL, `site_name` text NOT NULL, `location` text NOT NULL, `coordinates` text NOT NULL);
