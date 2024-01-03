-- +goose Up
-- add column "name" to table: "dahua_device_credentials"
ALTER TABLE `dahua_device_credentials` ADD COLUMN `name` text NOT NULL;

-- +goose Down
-- reverse: add column "name" to table: "dahua_device_credentials"
ALTER TABLE `dahua_device_credentials` DROP COLUMN `name`;
