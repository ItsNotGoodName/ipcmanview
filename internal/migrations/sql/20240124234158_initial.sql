-- +goose Up
-- add column "disabled_at" to table: "dahua_devices"
ALTER TABLE `dahua_devices` ADD COLUMN `disabled_at` datetime NULL;
-- add column "disabled_at" to table: "users"
ALTER TABLE `users` ADD COLUMN `disabled_at` datetime NULL;
-- add column "disabled_at" to table: "groups"
ALTER TABLE `groups` ADD COLUMN `disabled_at` datetime NULL;

-- +goose Down
-- reverse: add column "disabled_at" to table: "groups"
ALTER TABLE `groups` DROP COLUMN `disabled_at`;
-- reverse: add column "disabled_at" to table: "users"
ALTER TABLE `users` DROP COLUMN `disabled_at`;
-- reverse: add column "disabled_at" to table: "dahua_devices"
ALTER TABLE `dahua_devices` DROP COLUMN `disabled_at`;
