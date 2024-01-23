-- +goose Up
-- create index "group_users_user_id_group_id" to table: "group_users"
CREATE UNIQUE INDEX `group_users_user_id_group_id` ON `group_users` (`user_id`, `group_id`);

-- +goose Down
-- reverse: create index "group_users_user_id_group_id" to table: "group_users"
DROP INDEX `group_users_user_id_group_id`;
