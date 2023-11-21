-- +goose Up
-- create "dahua_cameras" table
CREATE TABLE `dahua_cameras` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `name` text NOT NULL, `address` text NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `location` text NOT NULL, `created_at` datetime NOT NULL);
-- create index "dahua_cameras_name" to table: "dahua_cameras"
CREATE UNIQUE INDEX `dahua_cameras_name` ON `dahua_cameras` (`name`);
-- create index "dahua_cameras_address" to table: "dahua_cameras"
CREATE UNIQUE INDEX `dahua_cameras_address` ON `dahua_cameras` (`address`);

-- +goose Down
-- reverse: create index "dahua_cameras_address" to table: "dahua_cameras"
DROP INDEX `dahua_cameras_address`;
-- reverse: create index "dahua_cameras_name" to table: "dahua_cameras"
DROP INDEX `dahua_cameras_name`;
-- reverse: create "dahua_cameras" table
DROP TABLE `dahua_cameras`;
