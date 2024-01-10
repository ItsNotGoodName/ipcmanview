-- +goose Up
-- add column "default" to table: "dahua_streams"
DELETE FROM `dahua_streams`;
ALTER TABLE `dahua_streams` ADD COLUMN `default` boolean NOT NULL;

-- +goose Down
-- reverse: add column "default" to table: "dahua_streams"
ALTER TABLE `dahua_streams` DROP COLUMN `default`;
