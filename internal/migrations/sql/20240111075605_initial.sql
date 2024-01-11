-- +goose Up
-- create index "dahua_afero_files_file_id" to table: "dahua_afero_files"
CREATE UNIQUE INDEX `dahua_afero_files_file_id` ON `dahua_afero_files` (`file_id`);
-- create index "dahua_afero_files_email_attachment_id" to table: "dahua_afero_files"
CREATE UNIQUE INDEX `dahua_afero_files_email_attachment_id` ON `dahua_afero_files` (`email_attachment_id`);

-- +goose Down
-- reverse: create index "dahua_afero_files_email_attachment_id" to table: "dahua_afero_files"
DROP INDEX `dahua_afero_files_email_attachment_id`;
-- reverse: create index "dahua_afero_files_file_id" to table: "dahua_afero_files"
DROP INDEX `dahua_afero_files_file_id`;
