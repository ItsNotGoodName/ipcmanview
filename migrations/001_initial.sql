-- Write your migrate up statements here

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE SCHEMA dahua;

CREATE TABLE dahua.cameras (
  id SERIAL PRIMARY KEY,
  address TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
--  sequence TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

----------- increment camera's sequence when it is updated
-- 
-- CREATE FUNCTION dahua.update_cameras_counter() RETURNS TRIGGER AS $$
--   BEGIN
--     NEW.sequence = now();
--     RETURN NEW;
--   END;
-- $$ LANGUAGE plpgsql;
-- 
-- CREATE TRIGGER tr_dahua_cameras_update
-- BEFORE UPDATE ON dahua.cameras
-- FOR EACH ROW
-- EXECUTE PROCEDURE dahua.update_cameras_counter();

--------- notify when camera is updated

CREATE FUNCTION dahua.notify_cameras_updated() RETURNS trigger as $$
  BEGIN
    PERFORM pg_notify('dahua.cameras:updated', cast(OLD.id AS TEXT));
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_dahua_cameras_updated
AFTER UPDATE ON dahua.cameras
FOR EACH ROW
EXECUTE PROCEDURE dahua.notify_cameras_updated();

--------- notify when camera is deleted

CREATE FUNCTION dahua.notify_cameras_deleted() RETURNS trigger as $$
  BEGIN
    PERFORM pg_notify('dahua.cameras:deleted', cast(OLD.id AS TEXT));
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_dahua_cameras_deleted
AFTER DELETE ON dahua.cameras
FOR EACH ROW
EXECUTE FUNCTION dahua.notify_cameras_deleted();

--------- 

CREATE TABLE dahua.camera_details (
  id INTEGER REFERENCES dahua.cameras(id) ON DELETE CASCADE,
  sn TEXT DEFAULT '' NOT NULL,
  device_class TEXT DEFAULT '' NOT NULL,
  device_type TEXT DEFAULT '' NOT NULL,
  hardware_version TEXT DEFAULT '' NOT NULL,
  market_area TEXT DEFAULT '' NOT NULL,
  process_info TEXT DEFAULT '' NOT NULL,
  vendor TEXT DEFAULT '' NOT NULL
);

--------- 

CREATE TABLE dahua.scanner(
  id INTEGER REFERENCES dahua.cameras(id) ON DELETE CASCADE,
  lock BOOLEAN GENERATED ALWAYS AS (lock_pid <> '') STORED,
  lock_pid TEXT NOT NULL,
  full_complete BOOLEAN DEFAULT false,
  full_cursor TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  quick_cursor TIMESTAMPTZ NOT NULL
);

CREATE TYPE dahua.scan_kind AS ENUM ('full', 'manual');

-- TODO: whenever a scan is completed, place the range cursor into dahua.scanner
CREATE TABLE dahua.completed_scans (
  id SERIAL PRIMARY KEY,
  camera_id INTEGER REFERENCES dahua.cameras(id) ON DELETE CASCADE,
  kind dahua.scan_kind NOT NULL,
  range tstzrange NOT NULL,
  range_cursor TIMESTAMPTZ NOT NULL,
  started_at TIMESTAMPTZ NOT NULL,
  deleted INTEGER NOT NULL,
  upserted INTEGER NOT NULL,
  percent REAL NOT NULL,
  duration INTEGER NOT NULL,
  success BOOLEAN GENERATED ALWAYS AS (error = '') STORED,
  error TEXT NOT NULL
);


---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
