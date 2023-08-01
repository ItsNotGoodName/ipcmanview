-- Write your migrate up statements here

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE dahua_cameras (
  id SERIAL PRIMARY KEY,
  address TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE FUNCTION notify_dahua_cameras_deleted() RETURNS trigger as $$
  BEGIN
    perform pg_notify('dahua_cameras:deleted', cast(OLD.id AS text));
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER dahua_cameras_deleted
AFTER DELETE ON dahua_cameras
FOR EACH ROW
EXECUTE FUNCTION notify_dahua_cameras_deleted();

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
