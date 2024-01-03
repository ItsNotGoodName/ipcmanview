-- +goose Up
-- create index "dahua_credentials_name" to table: "dahua_credentials"
CREATE UNIQUE INDEX `dahua_credentials_name` ON `dahua_credentials` (`name`);
-- create index "dahua_device_credentials_name" to table: "dahua_device_credentials"
CREATE UNIQUE INDEX `dahua_device_credentials_name` ON `dahua_device_credentials` (`name`);

-- +goose Down
-- reverse: create index "dahua_device_credentials_name" to table: "dahua_device_credentials"
DROP INDEX `dahua_device_credentials_name`;
-- reverse: create index "dahua_credentials_name" to table: "dahua_credentials"
DROP INDEX `dahua_credentials_name`;
