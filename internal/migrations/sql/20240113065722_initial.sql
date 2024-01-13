-- +goose Up
-- add column "ready" to table: "dahua_afero_files"
ALTER TABLE `dahua_afero_files` ADD COLUMN `ready` boolean NOT NULL;

-- +goose Down
-- reverse: add column "ready" to table: "dahua_afero_files"
ALTER TABLE `dahua_afero_files` DROP COLUMN `ready`;
