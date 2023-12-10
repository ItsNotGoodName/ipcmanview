CREATE TABLE settings (
  site_name TEXT NOT NULL,
  default_location TEXT NOT NULL
);

CREATE TABLE dahua_cameras (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  address TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  location TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE TABLE dahua_seeds (
  seed INTEGER NOT NULL PRIMARY KEY,
  camera_id INTEGER UNIQUE,
  FOREIGN KEY(camera_id) REFERENCES dahua_cameras(id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE dahua_events (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  camera_id INTEGER NOT NULL,
  code TEXT NOT NULL,
  action TEXT NOT NULL,
  `index` INTEGER NOT NULL,
  data JSON NOT NULL,
  created_at DATETIME NOT NULL,

  FOREIGN KEY(camera_id) REFERENCES dahua_cameras(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_files (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  camera_id INTEGER NOT NULL,
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
  work_dir_sn INTEGER NOT NULL,
  updated_at DATETIME NOT NULL,

  UNIQUE (camera_id, file_path),
  FOREIGN KEY(camera_id) REFERENCES dahua_cameras(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_file_cursors (
  camera_id INTEGER NOT NULL UNIQUE,
  quick_cursor DATETIME NOT NULL, -- (scanned) <- quick_cursor -> (not scanned / volatile)
  full_cursor DATETIME NOT NULL,  -- (not scanned) <- full_cursor -> (scanned)
  full_epoch DATETIME NOT NULL,
  full_complete BOOLEAN NOT NULL GENERATED ALWAYS AS (full_cursor <= full_epoch) STORED,

  FOREIGN KEY(camera_id) REFERENCES dahua_cameras(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_file_scan_locks (
  camera_id INTEGER NOT NULL UNIQUE,
  created_at DATETIME NOT NULL
);
