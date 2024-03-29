package models

import (
	"encoding/json"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager/config"
)

const (
	DahuaFileType_JPG = "jpg"
	DahuaFileType_DAV = "dav"
)

type DahuaFileSource string

const (
	DahuaFileSource_Device DahuaFileSource = "device"
	DahuaFileSource_Email  DahuaFileSource = "email"
)

type DahuaWorkerType string

const (
	DahuaWorkerType_Event     DahuaWorkerType = "event"
	DahuaWorkerType_Coaxial   DahuaWorkerType = "coaxial"
	DahuaWorkerType_QuickScan DahuaWorkerType = "quick-scan"
)

type DahuaWorkerState string

const (
	DahuaWorkerState_Connecting   DahuaWorkerState = "connecting"
	DahuaWorkerState_Connected    DahuaWorkerState = "connected"
	DahuaWorkerState_Disconnected DahuaWorkerState = "disconnected"
)

type DahuaFeature int

func (f DahuaFeature) EQ(feature DahuaFeature) bool {
	return feature != 0 && f&feature == feature
}

const (
	// DahuaFeature_Camera means the device is a camera.
	DahuaFeature_Camera DahuaFeature = 1 << iota
)

type DahuaScanType string

var (
	DahuaScanType_Unknown DahuaScanType = ""
	DahuaScanType_Full    DahuaScanType = "full"
	DahuaScanType_Quick   DahuaScanType = "quick"
	DahuaScanType_Reverse DahuaScanType = "reverse"
)

type DahuaPermissionLevel int

const (
	DahuaPermissionLevel_User DahuaPermissionLevel = iota
	DahuaPermissionLevel_Operator
	DahuaPermissionLevel_Admin
)

func (l DahuaPermissionLevel) String() string {
	switch l {
	case DahuaPermissionLevel_User:
		return "user"
	case DahuaPermissionLevel_Operator:
		return "operator"
	case DahuaPermissionLevel_Admin:
		return "admin"
	default:
		return "unknown"
	}
}

type DahuaError struct {
	Error string `json:"error"`
}

type DahuaRPCStatus struct {
	Error     string    `json:"error"`
	State     string    `json:"state"`
	LastLogin time.Time `json:"last_login"`
}

type DahuaDetail struct {
	SN               string `json:"sn"`
	DeviceClass      string `json:"device_class"`
	DeviceType       string `json:"device_type"`
	HardwareVersion  string `json:"hardware_version"`
	MarketArea       string `json:"market_area"`
	ProcessInfo      string `json:"process_info"`
	Vendor           string `json:"vendor"`
	OnvifVersion     string `json:"onvif_version"`
	AlgorithmVersion string `json:"algorithm_version"`
}

type DahuaSoftwareVersion struct {
	Build                   string `json:"build"`
	BuildDate               string `json:"build_date"`
	SecurityBaseLineVersion string `json:"security_base_line_version"`
	Version                 string `json:"version"`
	WebVersion              string `json:"web_version"`
}

type DahuaLicense struct {
	AbroadInfo    string    `json:"abroad_info"`
	AllType       bool      `json:"all_type"`
	DigitChannel  int       `json:"digit_channel"`
	EffectiveDays int       `json:"effective_days"`
	EffectiveTime time.Time `json:"effective_time"`
	LicenseID     int       `json:"license_id"`
	ProductType   string    `json:"product_type"`
	Status        int       `json:"status"`
	Username      string    `json:"username"`
}

type DahuaCoaxialStatus struct {
	WhiteLight bool `json:"white_light"`
	Speaker    bool `json:"speaker"`
}

type DahuaCoaxialCaps struct {
	SupportControlFullcolorLight bool `json:"support_control_fullcolor_light"`
	SupportControlLight          bool `json:"support_control_light"`
	SupportControlSpeaker        bool `json:"support_control_speaker"`
}

type DahuaFile struct {
	Channel     int       `json:"channel"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Length      int       `json:"length"`
	Type        string    `json:"type"`
	FilePath    string    `json:"file_path"`
	Duration    int       `json:"duration"`
	Disk        int       `json:"disk"`
	VideoStream string    `json:"video_stream"`
	Flags       []string  `json:"flags"`
	Events      []string  `json:"events"`
	Cluster     int       `json:"cluster"`
	Partition   int       `json:"partition"`
	PicIndex    int       `json:"pic_index"`
	Repeat      int       `json:"repeat"`
	WorkDir     string    `json:"work_dir"`
	WorkDirSN   bool      `json:"work_dir_sn"`
	Storage     Storage   `json:"storage"`
}

type DahuaStorage struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	Path       string `json:"path"`
	Type       string `json:"type"`
	TotalBytes int64  `json:"total_bytes"`
	UsedBytes  int64  `json:"used_bytes"`
	IsError    bool   `json:"is_error"`
}

type DahuaUser struct {
	ClientAddress string    `json:"client_address"`
	ClientType    string    `json:"client_type"`
	Group         string    `json:"group"`
	ID            int       `json:"id"`
	LoginTime     time.Time `json:"login_time"`
	Name          string    `json:"name"`
}

type DahuaSunriseSunset struct {
	SwitchMode  config.SwitchMode    `json:"switch_mode"`
	TimeSection dahuarpc.TimeSection `json:"time_section"`
}

type DahuaPreset struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

type DahuaEvent struct {
	ID        int64           `json:"id"`
	DeviceID  int64           `json:"device_id"`
	Code      string          `json:"code"`
	Action    string          `json:"action"`
	Index     int64           `json:"index"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
}

type DahuaUptime struct {
	Last      time.Time `json:"last"`
	Total     time.Time `json:"total"`
	Supported bool      `json:"supported"`
}
