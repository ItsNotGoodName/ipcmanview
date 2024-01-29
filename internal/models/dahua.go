package models

// TODO: remove all json tags

const DahuaFileTypeJPG = "jpg"
const DahuaFileTypeDAV = "dav"

type DahuaEventWorkerState string

const (
	DahuaEventWorkerStateConnecting   DahuaEventWorkerState = "connecting"
	DahuaEventWorkerStateConnected    DahuaEventWorkerState = "connected"
	DahuaEventWorkerStateDisconnected DahuaEventWorkerState = "disconnected"
)

type DahuaFeature int

func (f DahuaFeature) EQ(feature DahuaFeature) bool {
	return feature != 0 && f&feature == feature
}

const (
	// DahuaFeatureCamera means the device is a camera.
	DahuaFeatureCamera DahuaFeature = 1 << iota
)

type DahuaScanType string

var (
	DahuaScanTypeUnknown DahuaScanType = ""
	DahuaScanTypeFull    DahuaScanType = "full"
	DahuaScanTypeQuick   DahuaScanType = "quick"
	DahuaScanTypeReverse DahuaScanType = "reverse"
)

type DahuaPermissionLevel int

const (
	DahuaPermissionLevelUser DahuaPermissionLevel = iota
	DahuaPermissionLevelOperator
	DahuaPermissionLevelAdmin
)

func (l DahuaPermissionLevel) String() string {
	switch l {
	case DahuaPermissionLevelUser:
		return "user"
	case DahuaPermissionLevelOperator:
		return "operator"
	case DahuaPermissionLevelAdmin:
		return "admin"
	default:
		return "unknown"
	}
}
