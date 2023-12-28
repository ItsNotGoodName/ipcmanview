WITH RECURSIVE generate_series(value) AS (
  SELECT 1
  UNION ALL
  SELECT value+1 FROM generate_series WHERE value+1<=999
)
INSERT OR IGNORE INTO dahua_seeds (seed) SELECT value from generate_series;
INSERT OR IGNORE INTO dahua_event_rules (code) VALUES ('');
