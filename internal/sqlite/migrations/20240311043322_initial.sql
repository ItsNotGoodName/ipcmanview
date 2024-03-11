-- +goose Up
-- add column "email" to table: "dahua_devices"
ALTER TABLE `dahua_devices` ADD COLUMN `email` text NULL;
-- create index "dahua_devices_email" to table: "dahua_devices"
CREATE UNIQUE INDEX `dahua_devices_email` ON `dahua_devices` (`email`);

-- +goose Down
-- reverse: create index "dahua_devices_email" to table: "dahua_devices"
DROP INDEX `dahua_devices_email`;
-- reverse: add column "email" to table: "dahua_devices"
ALTER TABLE `dahua_devices` DROP COLUMN `email`;
