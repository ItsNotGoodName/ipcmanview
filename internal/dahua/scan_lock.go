package dahua

import "time"

// ScanaLock prevents concurrent scans on the same camera.
type ScanLock struct {
	CameraID  int64
	UUID      string
	CreatedAt time.Time
}
