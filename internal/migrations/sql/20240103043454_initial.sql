-- +goose Up
-- add column "name" to table: "dahua_credentials"
ALTER TABLE `dahua_credentials` ADD COLUMN `name` text NOT NULL;

-- +goose Down
-- reverse: add column "name" to table: "dahua_credentials"
ALTER TABLE `dahua_credentials` DROP COLUMN `name`;
