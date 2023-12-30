-- +goose Up
-- add column "name" to table: "dahua_streams"
ALTER TABLE `dahua_streams` ADD COLUMN `name` text NOT NULL DEFAULT "";

-- +goose Down
-- reverse: add column "name" to table: "dahua_streams"
ALTER TABLE `dahua_streams` DROP COLUMN `name`;
