-- +goose Up
-- create "squeuel" table
CREATE TABLE `squeuel` (`id` text NULL, `task_id` text NULL, `queue` text NOT NULL, `payload` blob NOT NULL, `timeout` datetime NOT NULL, `received` integer NOT NULL, `max_received` integer NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, PRIMARY KEY (`id`));
-- create index "squeuel_task_id" to table: "squeuel"
CREATE UNIQUE INDEX `squeuel_task_id` ON `squeuel` (`task_id`);
-- create index "squeuel_queue_created_idx" to table: "squeuel"
CREATE INDEX `squeuel_queue_created_idx` ON `squeuel` (`queue`, `created_at`);

-- +goose Down
-- reverse: create index "squeuel_queue_created_idx" to table: "squeuel"
DROP INDEX `squeuel_queue_created_idx`;
-- reverse: create index "squeuel_task_id" to table: "squeuel"
DROP INDEX `squeuel_task_id`;
-- reverse: create "squeuel" table
DROP TABLE `squeuel`;
