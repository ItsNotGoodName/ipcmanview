-- +goose Up
-- create index "dahua_files_device_id_start_time_idx" to table: "dahua_files"
CREATE INDEX `dahua_files_device_id_start_time_idx` ON `dahua_files` (`device_id`, `start_time`);

-- +goose Down
-- reverse: create index "dahua_files_device_id_start_time_idx" to table: "dahua_files"
DROP INDEX `dahua_files_device_id_start_time_idx`;
