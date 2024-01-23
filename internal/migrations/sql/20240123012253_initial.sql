-- +goose Up
-- create index "groups_name" to table: "groups"
CREATE UNIQUE INDEX `groups_name` ON `groups` (`name`);

-- +goose Down
-- reverse: create index "groups_name" to table: "groups"
DROP INDEX `groups_name`;
