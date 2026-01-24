package dahua

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/coaxialcontrolio"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/intervideo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/license"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/magicbox"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/peripheralchip"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/ptz"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/storage"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/usermanager"
)

type Feature string

const (
	// Device has at least 1 camera.
	FeatureCamera Feature = "camera"
)

func NewConn(v Device) Conn {
	return Conn{
		Key:      v.Key,
		Name:     v.Name,
		IP:       v.IP,
		Username: v.Username,
		Password: v.Password,
	}
}

type Conn struct {
	types.Key
	Name     string
	IP       string
	Username string
	Password string
}

func (lhs Conn) EQ(rhs Conn) bool {
	return lhs.Name == rhs.Name &&
		lhs.IP == rhs.IP &&
		lhs.Username == rhs.Username &&
		lhs.Password == rhs.Password
}

func isFatalError(err error) bool {
	res := &dahuarpc.Error{}
	if errors.As(err, &res) && slices.Contains([]dahuarpc.ErrorType{
		dahuarpc.ErrorTypeInvalidRequest,
		dahuarpc.ErrorTypeMethodNotFound,
		dahuarpc.ErrorTypeInterfaceNotFound,
		dahuarpc.ErrorTypeUnknown,
	}, res.Type) {
		slog.Error("Ignoring RPC ResponseError", slog.String("method", res.Method), slog.Int("code", res.Code), slog.String("type", string(res.Type)))
		return false
	}

	return true
}

func checkFatalError(err error) error {
	if isFatalError(err) {
		return err
	}
	return nil
}

type DeviceDetail struct {
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

func GetDetail(ctx context.Context, rpcClient dahuarpc.Conn) (DeviceDetail, error) {
	sn, err := magicbox.GetSerialNo(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return DeviceDetail{}, err
	}

	deviceClass, err := magicbox.GetDeviceClass(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return DeviceDetail{}, err
	}

	deviceType, err := magicbox.GetDeviceType(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return DeviceDetail{}, err
	}

	hardwareVersion, err := magicbox.GetHardwareVersion(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return DeviceDetail{}, err
	}

	marketArea, err := magicbox.GetMarketArea(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return DeviceDetail{}, err
	}

	processInfo, err := magicbox.GetProcessInfo(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return DeviceDetail{}, err
	}

	vendor, err := magicbox.GetVendor(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return DeviceDetail{}, err
	}

	onvifVersion, err := intervideo.ManagerGetVersion(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return DeviceDetail{}, err
	}

	var algorithmVersion string
	{
		res, err := peripheralchip.GetVersion(ctx, rpcClient, peripheralchip.TypeBLOB)
		if err != nil && isFatalError(err) {
			return DeviceDetail{}, err
		}
		if len(res) > 0 {
			algorithmVersion = res[0].SoftwareVersion
		}
	}

	return DeviceDetail{
		SN:               sn,
		DeviceClass:      deviceClass,
		DeviceType:       deviceType,
		HardwareVersion:  hardwareVersion,
		MarketArea:       marketArea,
		ProcessInfo:      processInfo,
		Vendor:           vendor,
		OnvifVersion:     onvifVersion,
		AlgorithmVersion: algorithmVersion,
	}, nil
}

type DeviceSoftwareVersion struct {
	Build                   string `json:"build"`
	BuildDate               string `json:"build_date"`
	SecurityBaseLineVersion string `json:"security_base_line_version"`
	Version                 string `json:"version"`
	WebVersion              string `json:"web_version"`
}

func GetSoftwareVersion(ctx context.Context, rpcClient dahuarpc.Conn) (DeviceSoftwareVersion, error) {
	res, err := magicbox.GetSoftwareVersion(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return DeviceSoftwareVersion{}, err
	}

	return DeviceSoftwareVersion{
		Build:                   res.Build,
		BuildDate:               res.BuildDate,
		SecurityBaseLineVersion: res.SecurityBaseLineVersion,
		Version:                 res.Version,
		WebVersion:              res.WebVersion,
	}, nil
}

type DeviceLicense struct {
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

func ListLicenses(ctx context.Context, rpcClient dahuarpc.Conn) ([]DeviceLicense, error) {
	licenses, err := license.GetLicenseInfo(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return nil, err
	}

	res := make([]DeviceLicense, 0, len(licenses))
	for _, l := range licenses {
		effectiveTime := time.Unix(int64(l.EffectiveTime), 0)

		res = append(res, DeviceLicense{
			AbroadInfo:    l.AbroadInfo,
			AllType:       l.AllType,
			DigitChannel:  l.DigitChannel,
			EffectiveDays: l.EffectiveDays,
			EffectiveTime: effectiveTime,
			LicenseID:     l.LicenseID,
			ProductType:   l.ProductType,
			Status:        l.Status,
			Username:      l.Username,
		})
	}

	return res, nil
}

type DeviceStorage struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	Path       string `json:"path"`
	Type       string `json:"type"`
	TotalBytes int64  `json:"total_bytes"`
	UsedBytes  int64  `json:"used_bytes"`
	IsError    bool   `json:"is_error"`
}

func ListStorage(ctx context.Context, rpcClient dahuarpc.Conn) ([]DeviceStorage, error) {
	devices, err := storage.GetDeviceAllInfo(ctx, rpcClient)
	if err != nil {
		return []DeviceStorage{}, checkFatalError(err)
	}

	res := []DeviceStorage{}
	for _, device := range devices {
		for _, detail := range device.Detail {
			res = append(res, DeviceStorage{
				Name:       device.Name,
				State:      device.State,
				Path:       detail.Path,
				Type:       detail.Type,
				TotalBytes: detail.TotalBytes.Integer(),
				UsedBytes:  detail.UsedBytes.Integer(),
				IsError:    detail.IsError,
			})
		}
	}

	return res, nil
}

type DeviceCoaxialStatus struct {
	WhiteLight bool `json:"white_light"`
	Speaker    bool `json:"speaker"`
}

func GetCoaxialStatus(ctx context.Context, rpcClient dahuarpc.Conn, channel int) (DeviceCoaxialStatus, error) {
	status, err := coaxialcontrolio.GetStatus(ctx, rpcClient, channel)
	if err != nil && isFatalError(err) {
		return DeviceCoaxialStatus{}, err
	}

	return DeviceCoaxialStatus{
		Speaker:    status.Speaker == "On",
		WhiteLight: status.WhiteLight == "On",
	}, nil
}

type DeviceCoaxialCaps struct {
	SupportControlFullcolorLight bool `json:"support_control_fullcolor_light"`
	SupportControlLight          bool `json:"support_control_light"`
	SupportControlSpeaker        bool `json:"support_control_speaker"`
}

func GetCoaxialCaps(ctx context.Context, rpcClient dahuarpc.Conn, channel int) (DeviceCoaxialCaps, error) {
	caps, err := coaxialcontrolio.GetCaps(ctx, rpcClient, channel)
	if err != nil && isFatalError(err) {
		return DeviceCoaxialCaps{}, err
	}

	return DeviceCoaxialCaps{
		SupportControlFullcolorLight: caps.SupportControlFullcolorLight == 1,
		SupportControlLight:          caps.SupportControlLight == 1,
		SupportControlSpeaker:        caps.SupportControlSpeaker == 1,
	}, nil
}

type DeviceActiveUser struct {
	ClientAddress string    `json:"client_address"`
	ClientType    string    `json:"client_type"`
	Group         string    `json:"group"`
	ID            int       `json:"id"`
	LoginTime     time.Time `json:"login_time"`
	Name          string    `json:"name"`
}

func ListActiveUsers(ctx context.Context, rpcClient dahuarpc.Conn, location *time.Location) ([]DeviceActiveUser, error) {
	users, err := usermanager.GetActiveUserInfoAll(ctx, rpcClient)
	if err != nil {
		return []DeviceActiveUser{}, checkFatalError(err)
	}

	res := make([]DeviceActiveUser, 0, len(users))
	for _, u := range users {
		loginTime, err := u.LoginTime.Parse(location)
		if err != nil {
			return nil, err
		}

		res = append(res, DeviceActiveUser{
			ClientAddress: u.ClientAddress,
			ClientType:    u.ClientType,
			Group:         u.Group,
			ID:            u.ID,
			LoginTime:     loginTime,
			Name:          u.Name,
		})
	}

	return res, nil
}

type DeviceUser struct {
	Anonymous            bool      `json:"anonymous"`
	AuthorityList        []string  `json:"authority_list"`
	Group                string    `json:"group"`
	ID                   int       `json:"id"`
	Memo                 string    `json:"memo"`
	Name                 string    `json:"name"`
	Password             string    `json:"password"`
	PasswordModifiedTime time.Time `json:"password_modified_time"`
	PwdScore             int       `json:"pwd_score"`
	Reserved             bool      `json:"reserved"`
	Sharable             bool      `json:"sharable"`
}

func ListUsers(ctx context.Context, rpcClient dahuarpc.Conn, location *time.Location) ([]DeviceUser, error) {
	users, err := usermanager.GetUserInfoAll(ctx, rpcClient)
	if err != nil {
		return []DeviceUser{}, checkFatalError(err)
	}

	res := make([]DeviceUser, 0, len(users))
	for _, u := range users {
		passwordModifiedTime, err := u.PasswordModifiedTime.Parse(location)
		if err != nil {
			return nil, err
		}
		res = append(res, DeviceUser{
			Anonymous:            u.Anonymous,
			AuthorityList:        u.AuthorityList,
			Group:                u.Group,
			ID:                   u.ID,
			Memo:                 u.Memo,
			Name:                 u.Name,
			Password:             u.Password,
			PasswordModifiedTime: passwordModifiedTime,
			PwdScore:             u.PwdScore,
			Reserved:             u.Reserved,
			Sharable:             u.Sharable,
		})
	}

	return res, nil
}

type DeviceGroup struct {
	AuthorityList []string `json:"authority_list"`
	ID            int      `json:"id"`
	Memo          string   `json:"memo"`
	Name          string   `json:"name"`
}

func ListGroups(ctx context.Context, rpcClient dahuarpc.Conn) ([]DeviceGroup, error) {
	users, err := usermanager.GetGroupInfoAll(ctx, rpcClient)
	if err != nil {
		return []DeviceGroup{}, checkFatalError(err)
	}

	res := make([]DeviceGroup, 0, len(users))
	for _, u := range users {
		res = append(res, DeviceGroup{
			AuthorityList: u.AuthorityList,
			ID:            u.ID,
			Memo:          u.Memo,
			Name:          u.Name,
		})
	}

	return res, nil
}

type DeviceStatus struct {
	Error     string    `json:"error"`
	State     string    `json:"state"`
	LastLogin time.Time `json:"last_login"`
}

func GetStatus(ctx context.Context, rpcClient dahuarpc.Client) DeviceStatus {
	rpcState := rpcClient.State(ctx)
	var rpcError string
	if rpcState.Error != nil {
		rpcError = rpcState.Error.Error()
	}
	return DeviceStatus{
		Error:     rpcError,
		State:     rpcState.State.String(),
		LastLogin: rpcState.LastLogin,
	}
}

type DevicePTZPreset struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

func ListPTZPresets(ctx context.Context, clientPTZ ptz.Client, channel int) ([]DevicePTZPreset, error) {
	vv, err := ptz.GetPresets(ctx, clientPTZ, channel)
	if err != nil {
		return []DevicePTZPreset{}, checkFatalError(err)
	}
	res := make([]DevicePTZPreset, 0, len(vv))
	for _, v := range vv {
		res = append(res, DevicePTZPreset{
			Index: v.Index,
			Name:  v.Name,
		})
	}
	return res, nil
}

func SetPTZPreset(ctx context.Context, clientPTZ ptz.Client, channel, index int) error {
	return ptz.Start(ctx, clientPTZ, channel, ptz.Params{
		Code: "GotoPreset",
		Arg1: index,
	})
}

type DeviceUptime struct {
	Last      time.Time `json:"last"`
	Total     time.Time `json:"total"`
	Supported bool      `json:"supported"`
}

func GetUptime(ctx context.Context, c dahuarpc.Conn) (DeviceUptime, error) {
	uptime, err := magicbox.GetUpTime(ctx, c)
	if err != nil {
		return DeviceUptime{}, checkFatalError(err)
	}

	now := time.Now()

	return DeviceUptime{
		Last:      now.Add(-time.Duration(uptime.Last) * time.Second),
		Total:     now.Add(-time.Duration(uptime.Total) * time.Second),
		Supported: true,
	}, nil
}

func RebootDevice(ctx context.Context, c dahuarpc.Conn) error {
	_, err := magicbox.Reboot(ctx, c)
	return err
}

type Storage string

const (
	StorageLocal Storage = "local"
	StorageSFTP  Storage = "sftp"
	StorageFTP   Storage = "ftp"
	StorageNFS   Storage = "nfs"
	StorageSMB   Storage = "smb"
)

func StorageFromFilePath(filePath string) Storage {
	if strings.HasPrefix(filePath, "sftp://") {
		return StorageSFTP
	}
	if strings.HasPrefix(filePath, "ftp://") {
		return StorageFTP
	}
	if strings.HasPrefix(filePath, "nfs://") {
		return StorageNFS
	}
	if strings.HasPrefix(filePath, "smb://") {
		return StorageSMB
	}
	return StorageLocal
}

func ListEventCodes(ctx context.Context, db *sqlx.DB) ([]string, error) {
	var codes []string
	err := db.Select(&codes, `
		SELECT DISTINCT code FROM dahua_events
	`)
	return codes, err
}

func ListEventActions(ctx context.Context, db *sqlx.DB) ([]string, error) {
	var actions []string
	err := db.Select(&actions, `
		SELECT DISTINCT action FROM dahua_events
	`)
	return actions, err
}
