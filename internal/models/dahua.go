package models

import (
	"encoding/json"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

type DahuaStatus struct {
	CameraID     int64     `json:"camera_id"`
	Address      string    `json:"address"`
	Username     string    `json:"username"`
	Location     string    `json:"location"`
	Seed         int       `json:"seed"`
	CreatedAt    time.Time `json:"created_at"`
	RPCError     string    `json:"rpc_error"`
	RPCState     string    `json:"rpc_state"`
	RPCLastLogin time.Time `json:"rpc_last_login"`
}

type DahuaCamera struct {
	ID        int64
	Address   string `validate:"address"`
	Username  string
	Password  string
	Location  types.Location
	Seed      int
	CreatedAt time.Time
}

type DahuaDetail struct {
	CameraID         int64  `json:"camera_id"`
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
	CameraID                int64  `json:"camera_id"`
	Build                   string `json:"build"`
	BuildDate               string `json:"build_date"`
	SecurityBaseLineVersion string `json:"security_base_line_version"`
	Version                 string `json:"version"`
	WebVersion              string `json:"web_version"`
}

type DahuaLicense struct {
	CameraID      int64     `json:"camera_id"`
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
	CameraID   int64 `json:"camera_id"`
	WhiteLight bool  `json:"white_light"`
	Speaker    bool  `json:"speaker"`
}

type DahuaCoaxialCaps struct {
	CameraID                     int64 `json:"camera_id"`
	SupportControlFullcolorLight bool  `json:"support_control_fullcolor_light"`
	SupportControlLight          bool  `json:"support_control_light"`
	SupportControlSpeaker        bool  `json:"support_control_speaker"`
}

type DahuaFile struct {
	CameraID    int64     `json:"camera_id"`
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
	WorkDirSN   int       `json:"work_dir_sn"`
}

type DahuaEvent struct {
	ID            int64           `json:"id"`
	CameraID      int64           `json:"camera_id"`
	ContentType   string          `json:"content_type"`
	ContentLength int             `json:"content_length"`
	Code          string          `json:"code"`
	Action        string          `json:"action"`
	Index         int             `json:"index"`
	Data          json.RawMessage `json:"data"`
	CreatedAt     time.Time       `json:"created_at"`
}

type DahuaStorage struct {
	CameraID   int64  `json:"camera_id"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Path       string `json:"path"`
	Type       string `json:"type"`
	TotalBytes int64  `json:"total_bytes"`
	UsedBytes  int64  `json:"used_bytes"`
	IsError    bool   `json:"is_error"`
}

type DahuaUser struct {
	CameraID      int64     `json:"camera_id"`
	ClientAddress string    `json:"client_address"`
	ClientType    string    `json:"client_type"`
	Group         string    `json:"group"`
	ID            int       `json:"id"`
	LoginTime     time.Time `json:"login_time"`
	Name          string    `json:"name"`
}
