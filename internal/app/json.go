package app

import (
	"encoding/json"
	"net"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/build"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/coaxialcontrolio"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
)

type ManualFileScan struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}

type DeviceEvent struct {
	ID         string          `json:"id"`
	DeviceUUID string          `json:"device_uuid"`
	Code       string          `json:"code"`
	Action     string          `json:"action"`
	Index      int64           `json:"index"`
	Data       json.RawMessage `json:"data"`
	CreatedAt  time.Time       `json:"created_at"`
}

type DeviceVideoInModeSync struct {
	Location      *types.Location `json:"location,omitempty"`
	Latitude      *float64        `json:"latitude,omitempty"`
	Longitude     *float64        `json:"longitude,omitempty"`
	SunriseOffset *types.Duration `json:"sunrise_offset,omitempty"`
	SunsetOffset  *types.Duration `json:"sunset_offset,omitempty"`
}

type DeviceCoaxialControlInput struct {
	Channel  int                           `json:"int"`
	Controls []DeviceCoaxialControlRequest `json:"controls"`
}

type DeviceCoaxialControlRequest struct {
	Type        coaxialcontrolio.Type        `json:"type"`
	IO          coaxialcontrolio.IO          `json:"io"`
	TriggerMode coaxialcontrolio.TriggerMode `json:"trigger_mode"`
}

type FileScanCursor struct {
	QuickCursor  time.Time `json:"quick_cursor"`
	FullCursor   time.Time `json:"full_cursor"`
	FullEpoch    time.Time `json:"full_epoch"`
	FullComplete bool      `json:"full_complete"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type File struct{}

//---------- CRUD

func NewPagePagination(result pagination.PageResult) PagePagination {
	return PagePagination{
		Page:         result.Page,
		PerPage:      result.PerPage,
		TotalPages:   result.TotalPages,
		TotalItems:   result.TotalItems,
		SeenItems:    result.Seen(),
		PreviousPage: result.Previous(),
		NextPage:     result.Next(),
	}
}

type PagePagination struct {
	Page         int `json:"page"`
	PerPage      int `json:"per_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	SeenItems    int `json:"seen_items"`
	PreviousPage int `json:"previous_page"`
	NextPage     int `json:"next_page"`
}

type CreateDevice struct {
	UUID            *string         `json:"uuid,omitempty" format:"uuid"`
	Name            string          `json:"name,omitempty"`
	IP              string          `json:"ip,omitempty" format:"ipv4"`
	Username        string          `json:"username,omitempty"`
	Password        string          `json:"password,omitempty"`
	Location        *types.Location `json:"location,omitempty"`
	Features        []dahua.Feature `json:"features,omitempty" enum:"camera"`
	Email           *string         `json:"email,omitempty"`
	Latitude        *float64        `json:"latitude,omitempty"`
	Longitude       *float64        `json:"longitude,omitempty"`
	SunriseOffset   *types.Duration `json:"sunrise_offset,omitempty"`
	SunsetOffset    *types.Duration `json:"sunset_offset,omitempty"`
	SyncVideoInMode *bool           `json:"sync_video_in_mode,omitempty"`
}

type Device struct {
	UUID            string          `json:"uuid"`
	Name            string          `json:"name"`
	IP              net.IP          `json:"ip"`
	Username        string          `json:"username"`
	Location        *types.Location `json:"location"`
	Features        []dahua.Feature `json:"features"`
	Email           string          `json:"email"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Latitude        *float64        `json:"latitude"`
	Longitude       *float64        `json:"longitude"`
	SunriseOffset   *types.Duration `json:"sunrise_offset"`
	SunsetOffset    *types.Duration `json:"sunset_offset"`
	SyncVideoInMode *bool           `json:"sync_video_in_mode"`
}

type UpdateDevice struct {
	Name            string          `json:"name"`
	IP              string          `json:"ip" format:"ipv4"`
	Username        string          `json:"username"`
	Password        *string         `json:"password,omitempty"`
	Location        *types.Location `json:"location,omitempty"`
	Features        []dahua.Feature `json:"features,omitempty" enum:"camera"`
	Email           *string         `json:"email,omitempty"`
	Latitude        *float64        `json:"latitude"`
	Longitude       *float64        `json:"longitude"`
	SunriseOffset   *types.Duration `json:"sunrise_offset,omitempty"`
	SunsetOffset    *types.Duration `json:"sunset_offset,omitempty"`
	SyncVideoInMode *bool           `json:"sync_video_in_mode,omitempty"`
}

type CreateEmailEndpoint struct {
	UUID          *string  `json:"uuid,omitempty" format:"uuid"`
	URLs          []string `json:"urls"`
	Expression    string   `json:"expression,omitempty"`
	TitleTemplate *string  `json:"title_template,omitempty"`
	BodyTemplate  *string  `json:"body_template,omitempty"`
	Attachments   bool     `json:"attachments,omitempty"`
	DeviceUUIDs   []string `json:"device_uuids,omitempty"`
	Global        bool     `json:"global,omitempty"`
	Disabled      bool     `json:"disabled,omitempty" default:"false"`
}

type EmailEndpoint struct {
	UUID          string    `json:"uuid"`
	Global        bool      `json:"global"`
	DeviceUUIDs   []string  `json:"device_uuids"`
	Expression    string    `json:"expression"`
	URLs          []string  `json:"urls"`
	TitleTemplate string    `json:"title_template"`
	BodyTemplate  string    `json:"body_template"`
	Attachments   bool      `json:"attachments"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Disabled      bool      `json:"disabled"`
}

type UpdateEmailEndpoint struct {
	URLs          []string `json:"urls"`
	Expression    string   `json:"expression,omitempty"`
	TitleTemplate string   `json:"title_template,omitempty"`
	BodyTemplate  string   `json:"body_template,omitempty"`
	Attachments   bool     `json:"attachments,omitempty"`
	DeviceUUIDs   []string `json:"device_uuids,omitempty"`
	Global        bool     `json:"global,omitempty"`
	Disabled      bool     `json:"disabled,omitempty" default:"false"`
}

type CreateStorageDestination struct {
	UUID            *string `json:"uuid,omitempty" format:"uuid"`
	Name            string  `json:"name"`
	Storage         string  `json:"storage" enum:"sftp,ftp"`
	ServerAddress   string  `json:"server_address"`
	Port            int     `json:"port"`
	Username        string  `json:"username"`
	Password        string  `json:"password"`
	RemoteDirectory string  `json:"remote_directory"`
}

type StorageDestination struct {
	UUID            string    `json:"uuid"`
	Name            string    `json:"name"`
	Storage         string    `json:"storage"`
	ServerAddress   string    `json:"server_address"`
	Port            int       `json:"port"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	RemoteDirectory string    `json:"remote_directory"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UpdateStorageDestination struct {
	Name            string `json:"name"`
	Storage         string `json:"storage" enum:"sftp,ftp"`
	ServerAddress   string `json:"server_address"`
	Port            int    `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	RemoteDirectory string `json:"remote_directory"`
}

type Settings struct {
	Location        types.Location `json:"location"`
	Latitude        float64        `json:"latitude"`
	Longitude       float64        `json:"longitude"`
	SunriseOffset   types.Duration `json:"sunrise_offset"`
	SunsetOffset    types.Duration `json:"sunset_offset"`
	SyncVideoInMode bool           `json:"sync_video_in_mode"`
}

type PatchSettings struct {
	Location        *types.Location `json:"location,omitempty"`
	Latitude        *float64        `json:"latitude,omitempty"`
	Longitude       *float64        `json:"longitude,omitempty"`
	SunriseOffset   *types.Duration `json:"sunrise_offset,omitempty"`
	SunsetOffset    *types.Duration `json:"sunset_offset,omitempty"`
	SyncVideoInMode *bool           `json:"sync_video_in_mode,omitempty"`
}

type UpdateSettings struct {
	Location        types.Location `json:"location"`
	Latitude        float64        `json:"latitude"`
	Longitude       float64        `json:"longitude"`
	SunriseOffset   types.Duration `json:"sunrise_offset"`
	SunsetOffset    types.Duration `json:"sunset_offset"`
	SyncVideoInMode bool           `json:"sync_video_in_mode"`
}

type CreateEventRule struct {
	UUID      *string `json:"uuid,omitempty" format:"uuid"`
	Code      string  `json:"code"`
	AllowDB   bool    `json:"allow_db"`
	AllowLive bool    `json:"allow_live"`
	AllowMQTT bool    `json:"allow_mqtt"`
}

type UpdateEventRule struct {
	Code      string `json:"code"`
	AllowDB   bool   `json:"allow_db"`
	AllowLive bool   `json:"allow_live"`
	AllowMQTT bool   `json:"allow_mqtt"`
}

//---------- Pages

type GetHomePage struct {
	Device_Count int               `json:"device_count"`
	Event_Count  int               `json:"event_count"`
	Email_Count  int               `json:"email_count"`
	File_Count   int               `json:"file_count"`
	DB_Usage     int               `json:"db_usage"`
	FileUsage    int64             `json:"file_usage"`
	Build        build.Build       `json:"build"`
	Emails       GetHomePageEmails `json:"home"`
}

type GetHomePageEmails struct{}
