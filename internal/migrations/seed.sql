-- +goose Up
WITH RECURSIVE generate_series(value) AS (
  SELECT 1
  UNION ALL
  SELECT value+1 FROM generate_series WHERE value+1<=999
)
INSERT INTO dahua_seeds (seed) SELECT value from generate_series;

-- +goose Down
WITH RECURSIVE generate_series(value) AS (
  SELECT 1
  UNION ALL
  SELECT value+1 FROM generate_series WHERE value+1<=999
)
DELETE FROM dahua_seeds WHERE seed IN (SELECT value from generate_series);
