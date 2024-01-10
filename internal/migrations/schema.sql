CREATE TABLE settings (
  site_name TEXT NOT NULL,
  default_location TEXT NOT NULL
);

CREATE TABLE dahua_devices (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  address TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  location TEXT NOT NULL,
  feature INTEGER NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE TABLE dahua_seeds (
  seed INTEGER NOT NULL PRIMARY KEY,
  device_id INTEGER UNIQUE,

  FOREIGN KEY(device_id) REFERENCES dahua_devices(id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE dahua_events (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  device_id INTEGER NOT NULL,
  code TEXT NOT NULL,
  action TEXT NOT NULL,
  `index` INTEGER NOT NULL,
  data JSON NOT NULL,
  created_at DATETIME NOT NULL,

  FOREIGN KEY(device_id) REFERENCES dahua_devices(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_event_rules(
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  code TEXT NOT NULL UNIQUE,
  ignore_db BOOLEAN NOT NULL DEFAULT false,
  ignore_live BOOLEAN NOT NULL DEFAULT false,
  ignore_mqtt BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE dahua_event_device_rules(
  device_id INTEGER NOT NULL,
  code TEXT NOT NULL,
  ignore_db BOOLEAN NOT NULL DEFAULT false,
  ignore_live BOOLEAN NOT NULL DEFAULT false,
  ignore_mqtt BOOLEAN NOT NULL DEFAULT false,

  UNIQUE (device_id, code),
  FOREIGN KEY(device_id) REFERENCES dahua_devices(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_event_worker_states (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  device_id INTEGER NOT NULL,
  state TEXT NOT NULL,
  error TEXT,
  created_at DATETIME NOT NULL,

  FOREIGN KEY(device_id) REFERENCES dahua_devices(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_files (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  device_id INTEGER NOT NULL,
  channel INTEGER NOT NULL,
  start_time DATETIME NOT NULL UNIQUE,
  end_time DATETIME NOT NULL,
  length INTEGER NOT NULL,
  type TEXT NOT NULL,
  file_path TEXT NOT NULL,
  duration INTEGER NOT NULL,
  disk INTEGER NOT NULL,
  video_stream TEXT NOT NULL,
  flags JSON NOT NULL,
  events JSON NOT NULL,
  cluster INTEGER NOT NULL,
  partition INTEGER NOT NULL,
  pic_index INTEGER NOT NULL,
  repeat INTEGER NOT NULL,
  work_dir TEXT NOT NULL,
  work_dir_sn BOOLEAN NOT NULL,
  updated_at DATETIME NOT NULL,
  storage TEXT NOT NULL,

  UNIQUE (device_id, file_path),
  FOREIGN KEY(device_id) REFERENCES dahua_devices(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_file_cursors (
  device_id INTEGER NOT NULL UNIQUE,
  quick_cursor DATETIME NOT NULL, -- (scanned) <- quick_cursor -> (not scanned / volatile)
  full_cursor DATETIME NOT NULL,  -- (not scanned) <- full_cursor -> (scanned)
  full_epoch DATETIME NOT NULL,
  full_complete BOOLEAN NOT NULL GENERATED ALWAYS AS (full_cursor <= full_epoch) STORED,

  scan BOOLEAN NOT NULL,
  scan_percent REAL NOT NULL,
  scan_type TEXT NOT NULL,

  FOREIGN KEY(device_id) REFERENCES dahua_devices(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_file_scan_locks (
  device_id INTEGER NOT NULL UNIQUE,
  touched_at DATETIME NOT NULL, 

  FOREIGN KEY(device_id) REFERENCES dahua_devices(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_storage_destinations (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  storage TEXT NOT NULL,
  server_address TEXT NOT NULL,
  port INTEGER NOT NULL,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  remote_directory TEXT NOT NULL
);

CREATE TABLE dahua_streams (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  internal BOOLEAN NOT NULL,
  device_id INTEGER NOT NULL,
  channel INTEGER NOT NULL,
  subtype INTEGER NOT NULL,
  name TEXT NOT NULL,
  mediamtx_path TEXT NOT NULL,

  UNIQUE(device_id, channel, subtype),
  FOREIGN KEY(device_id) REFERENCES dahua_devices(id) ON UPDATE CASCADE ON DELETE CASCADE
);
