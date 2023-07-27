-- Write your migrate up statements here

CREATE TABLE placeholder (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL DEFAULT ''
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
