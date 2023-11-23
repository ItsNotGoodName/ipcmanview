-- +goose Up
-- create index "dahua_files_start_time" to table: "dahua_files"
CREATE UNIQUE INDEX `dahua_files_start_time` ON `dahua_files` (`start_time`);

-- +goose Down
-- reverse: create index "dahua_files_start_time" to table: "dahua_files"
DROP INDEX `dahua_files_start_time`;
