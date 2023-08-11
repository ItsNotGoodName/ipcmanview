-- Write your migrate up statements here

CREATE TABLE users (
  id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  email TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE SCHEMA dahua;

CREATE TABLE dahua.cameras (
  id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  address TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  location TEXT NOT NULL, -- location of camera, stored as 'America/Los_Angeles' format and converted to Go's *time.Location, when this is updated, every camera files for this camera should be rescanned
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--------- notify when camera is updated

CREATE FUNCTION dahua.fn_cameras_updated() RETURNS TRIGGER AS $$
  BEGIN
    PERFORM pg_notify('dahua.cameras:updated', cast(OLD.id AS TEXT));
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_dahua_cameras_updated
AFTER UPDATE ON dahua.cameras
FOR EACH ROW
EXECUTE PROCEDURE dahua.fn_cameras_updated();

--------- notify when camera is deleted

CREATE FUNCTION dahua.fn_cameras_deleted() RETURNS TRIGGER AS $$
  BEGIN
    PERFORM pg_notify('dahua.cameras:deleted', cast(OLD.id AS TEXT));
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_dahua_cameras_deleted
AFTER DELETE ON dahua.cameras
FOR EACH ROW
EXECUTE FUNCTION dahua.fn_cameras_deleted();

--------- information cached from camera

CREATE TABLE dahua.camera_details (
  camera_id INTEGER NOT NULL UNIQUE REFERENCES dahua.cameras(id) ON DELETE CASCADE,
  sn TEXT DEFAULT '' NOT NULL,
  device_class TEXT DEFAULT '' NOT NULL,
  device_type TEXT DEFAULT '' NOT NULL,
  hardware_version TEXT DEFAULT '' NOT NULL,
  market_area TEXT DEFAULT '' NOT NULL,
  process_info TEXT DEFAULT '' NOT NULL,
  vendor TEXT DEFAULT '' NOT NULL
);

CREATE TABLE dahua.camera_softwares (
  camera_id INTEGER NOT NULL UNIQUE REFERENCES dahua.cameras(id) ON DELETE CASCADE,
  build TEXT DEFAULT '' NOT NULL,
  build_date TEXT DEFAULT '' NOT NULL,
  security_base_line_version TEXT DEFAULT '' NOT NULL,
  version TEXT DEFAULT '' NOT NULL,
  web_version TEXT DEFAULT '' NOT NULL
);

CREATE TABLE IF NOT EXISTS dahua.camera_licenses (
  camera_id INTEGER NOT NULL REFERENCES dahua.cameras(id) ON DELETE CASCADE,
  abroad_info TEXT DEFAULT '' NOT NULL,
  all_type BOOLEAN DEFAULT false NOT NULL,
  digit_channel INTEGER DEFAULT 0 NOT NULL,
  effective_days INTEGER DEFAULT 0 NOT NULL,
  effective_time INTEGER DEFAULT 0 NOT NULL,
  effective_at TIMESTAMPTZ NOT NULL GENERATED ALWAYS AS (TO_TIMESTAMP(effective_time)) STORED,
  license_id INTEGER DEFAULT 0 NOT NULL,
  product_type TEXT DEFAULT '' NOT NULL,
  status INTEGER DEFAULT 0 NOT NULL,
  username TEXT DEFAULT '' NOT NULL
);

--------- camera file scanning

CREATE TABLE dahua.scan_cursors (
  camera_id INTEGER NOT NULL UNIQUE REFERENCES dahua.cameras(id) ON DELETE CASCADE,
  quick_cursor TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP - INTERVAL '8 hours'),              -- (scanned) <- quick_cursor -> (not scanned / volatile)
  full_cursor TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP CHECK(full_cursor <= full_epoch_end), -- (not scanned) <- full_cursor -> (scanned)
  full_epoch TIMESTAMPTZ NOT NULL DEFAULT '2009-12-31 00:00:00',
  full_epoch_end TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  full_complete BOOLEAN NOT NULL GENERATED ALWAYS AS (full_cursor <= full_epoch) STORED
);

CREATE TABLE dahua.camera_files (
  id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  camera_id INTEGER NOT NULL REFERENCES dahua.scan_cursors(camera_id) ON DELETE CASCADE,
  file_path TEXT NOT NULL,
  kind TEXT NOT NULL,
  size INTEGER NOT NULL,
  start_time TIMESTAMPTZ NOT NULL UNIQUE,
  end_time TIMESTAMPTZ NOT NULL,
  duration INTEGER NOT NULL GENERATED ALWAYS AS (EXTRACT(EPOCH FROM end_time - start_time)) STORED,
  scanned_at TIMESTAMPTZ NOT NULL,
  events JSONB NOT NULL,
  UNIQUE (camera_id, file_path)
);


CREATE TABLE dahua.scan_seeds (
  seed INTEGER NOT NULL UNIQUE,
  camera_id INTEGER UNIQUE REFERENCES dahua.scan_cursors(camera_id) ON DELETE SET NULL
);

INSERT INTO dahua.scan_seeds (seed) VALUES (generate_series(1,999));

--------- notify when camera is created and insert 1-1 releations

CREATE FUNCTION dahua.fn_cameras_created() RETURNS TRIGGER AS $$
  BEGIN
    INSERT INTO dahua.camera_details (camera_id) VALUES (NEW.id);
    INSERT INTO dahua.camera_softwares (camera_id) VALUES (NEW.id);
    INSERT INTO dahua.scan_cursors (camera_id) VALUES (NEW.id);
    UPDATE dahua.scan_seeds SET camera_id = NEW.id 
      WHERE seed = (SELECT seed FROM dahua.scan_seeds WHERE camera_id = NEW.id OR camera_id IS NULL ORDER BY camera_id asc LIMIT 1);
    PERFORM pg_notify('dahua.cameras:created', cast(NEW.id AS TEXT));
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_dahua_cameras_created
AFTER INSERT ON dahua.cameras
FOR EACH ROW
EXECUTE PROCEDURE dahua.fn_cameras_created();

--------- 

CREATE TYPE dahua.scan_kind AS ENUM ('full', 'quick', 'manual');

CREATE TABLE dahua.scan_queue_tasks (
  id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  camera_id INTEGER NOT NULL REFERENCES dahua.scan_cursors(camera_id) ON DELETE CASCADE,
  kind dahua.scan_kind NOT NULL,
  range tstzrange NOT NULL DEFAULT 'empty'
);

CREATE TABLE dahua.scan_active_tasks (
  queue_id INTEGER NOT NULL REFERENCES dahua.scan_queue_tasks(id) ON DELETE CASCADE,
  camera_id INTEGER NOT NULL UNIQUE REFERENCES dahua.scan_cursors(camera_id) ON DELETE CASCADE,
  kind dahua.scan_kind NOT NULL,
  range tstzrange NOT NULL CHECK (range != 'empty'),
  cursor TIMESTAMPTZ NOT NULL,
  started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted INTEGER NOT NULL DEFAULT 0,
  upserted INTEGER NOT NULL DEFAULT 0,
  percent REAL NOT NULL DEFAULT 0.0
);

CREATE TABLE dahua.scan_complete_tasks (
  id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  camera_id INTEGER NOT NULL REFERENCES dahua.scan_cursors(camera_id) ON DELETE CASCADE,
  kind dahua.scan_kind NOT NULL,
  range tstzrange NOT NULL CHECK (range != 'empty'),
  cursor TIMESTAMPTZ NOT NULL,
  started_at TIMESTAMPTZ NOT NULL,
  deleted INTEGER NOT NULL,
  upserted INTEGER NOT NULL,
  percent REAL NOT NULL,
  duration INTEGER NOT NULL,
  success BOOLEAN GENERATED ALWAYS AS (error = '') STORED,
  error TEXT NOT NULL DEFAULT ''
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
