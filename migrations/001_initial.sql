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
  location TEXT NOT NULL, -- location of dahua camera, stored as 'America/Los_Angeles' format and converted to Go's *time.Location
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--------- notify when dahua camera is updated

CREATE FUNCTION dahua.notify_cameras_updated() RETURNS TRIGGER AS $$
  BEGIN
    PERFORM pg_notify('dahua.cameras:updated', cast(OLD.id AS TEXT));
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_dahua_cameras_updated
AFTER UPDATE ON dahua.cameras
FOR EACH ROW
EXECUTE PROCEDURE dahua.notify_cameras_updated();

--------- notify when dahua camera is deleted

CREATE FUNCTION dahua.notify_cameras_deleted() RETURNS TRIGGER AS $$
  BEGIN
    PERFORM pg_notify('dahua.cameras:deleted', cast(OLD.id AS TEXT));
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_dahua_cameras_deleted
AFTER DELETE ON dahua.cameras
FOR EACH ROW
EXECUTE FUNCTION dahua.notify_cameras_deleted();

--------- CACHED information about dahua cameras

CREATE TABLE dahua.camera_details (
  id INTEGER NOT NULL UNIQUE REFERENCES dahua.cameras(id) ON DELETE CASCADE,
  sn TEXT DEFAULT '' NOT NULL,
  device_class TEXT DEFAULT '' NOT NULL,
  device_type TEXT DEFAULT '' NOT NULL,
  hardware_version TEXT DEFAULT '' NOT NULL,
  market_area TEXT DEFAULT '' NOT NULL,
  process_info TEXT DEFAULT '' NOT NULL,
  vendor TEXT DEFAULT '' NOT NULL
);

CREATE TABLE dahua.camera_files (
  id SERIAL PRIMARY KEY,
  camera_id INTEGER NOT NULL REFERENCES dahua.cameras(id) ON DELETE CASCADE,
  file_path TEXT NOT NULL,
  kind TEXT NOT NULL,
  size INTEGER NOT NULL,
  start_time TIMESTAMPTZ NOT NULL UNIQUE,
  end_time TIMESTAMPTZ NOT NULL,
  duration INTEGER GENERATED ALWAYS AS (EXTRACT(EPOCH FROM start_time - end_time)) STORED,
  updated_at TIMESTAMPTZ NOT NULL,
  events JSONB NOT NULL,
  UNIQUE (camera_id, file_path)
);

--------- time seeds
-- Used to give dahua.camera_files a unique start_time.
CREATE TABLE dahua.time_seeds (
  seed INTEGER NOT NULL UNIQUE,
  camera_id INTEGER UNIQUE REFERENCES dahua.cameras(id) ON DELETE SET NULL
);
CREATE FUNCTION init_time_seeds() RETURNS void AS $$
BEGIN
  FOR i IN 1..999 LOOP
  INSERT INTO dahua.time_seeds (seed) VALUES (i);
  END LOOP;
END;
$$ LANGUAGE plpgsql;
SELECT init_time_seeds();
DROP FUNCTION init_time_seeds;

--------- dahua camera file scanning

CREATE TABLE dahua.scanners (
  id INTEGER NOT NULL UNIQUE REFERENCES dahua.cameras(id) ON DELETE CASCADE,
  full_complete BOOLEAN GENERATED ALWAYS AS (full_cursor <= full_epoch) STORED,
  full_cursor TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,                     -- (not scanned) <- full_cursor -> (scanned)
  full_epoch TIMESTAMPTZ NOT NULL DEFAULT '2009-12-31 00:00:00',
  quick_cursor TIMESTAMPTZ NOT NULL                                               -- (scanned) <- quick_cursor -> (not scanned / volatile)
);

CREATE TABLE dahua.scanner_locks (
  id INTEGER NOT NULL UNIQUE REFERENCES dahua.cameras(id) ON DELETE CASCADE,
  uuid UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE dahua.scan_kind AS ENUM ('full', 'manual');

CREATE TABLE dahua.completed_scans (
  id SERIAL PRIMARY KEY,
  camera_id INTEGER NOT NULL REFERENCES dahua.cameras(id) ON DELETE CASCADE,
  kind dahua.scan_kind NOT NULL,
  range tstzrange NOT NULL,
  -- range_cursor TIMESTAMPTZ NOT NULL,
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
