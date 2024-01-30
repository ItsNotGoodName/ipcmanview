package models

import (
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager/config"
)

type DahuaError struct {
	Error string `json:"error"`
}

type DahuaStatus struct {
	DeviceID     int64     `json:"device_id"`
	Url          string    `json:"url"`
	Username     string    `json:"username"`
	Location     string    `json:"location"`
	Seed         int       `json:"seed"`
	RPCError     string    `json:"rpc_error"`
	RPCState     string    `json:"rpc_state"`
	RPCLastLogin time.Time `json:"rpc_last_login"`
}

type DahuaDetail struct {
	DeviceID         int64  `json:"device_id"`
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
	DeviceID                int64  `json:"device_id"`
	Build                   string `json:"build"`
	BuildDate               string `json:"build_date"`
	SecurityBaseLineVersion string `json:"security_base_line_version"`
	Version                 string `json:"version"`
	WebVersion              string `json:"web_version"`
}

type DahuaLicense struct {
	DeviceID      int64     `json:"device_id"`
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
	DeviceID   int64 `json:"device_id"`
	WhiteLight bool  `json:"white_light"`
	Speaker    bool  `json:"speaker"`
}

type DahuaCoaxialCaps struct {
	DeviceID                     int64 `json:"device_id"`
	SupportControlFullcolorLight bool  `json:"support_control_fullcolor_light"`
	SupportControlLight          bool  `json:"support_control_light"`
	SupportControlSpeaker        bool  `json:"support_control_speaker"`
}

type DahuaFile struct {
	ID          int64     `json:"id"`
	DeviceID    int64     `json:"device_id"`
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
	DeviceID   int64  `json:"device_id"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Path       string `json:"path"`
	Type       string `json:"type"`
	TotalBytes int64  `json:"total_bytes"`
	UsedBytes  int64  `json:"used_bytes"`
	IsError    bool   `json:"is_error"`
}

type DahuaUser struct {
	DeviceID      int64     `json:"device_id"`
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