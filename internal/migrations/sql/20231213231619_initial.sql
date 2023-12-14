-- +goose Up
-- add column "touched_at" to table: "dahua_file_scan_locks"
ALTER TABLE `dahua_file_scan_locks` ADD COLUMN `touched_at` datetime NOT NULL;

-- +goose Down
-- reverse: add column "touched_at" to table: "dahua_file_scan_locks"
ALTER TABLE `dahua_file_scan_locks` DROP COLUMN `touched_at`;
