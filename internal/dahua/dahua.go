package dahua

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/coaxialcontrolio"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager/config"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/intervideo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/license"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/magicbox"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/peripheralchip"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/ptz"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/storage"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/usermanager"
	"github.com/nathan-osman/go-sunrise"
	"github.com/rs/zerolog/log"
)

func NewScanLockStore() ScanLockStore {
	return ScanLockStore{core.NewLockStore[int64]()}
}

type ScanLockStore struct{ *core.LockStore[int64] }

func isFatalError(err error) bool {
	res := &dahuarpc.ResponseError{}
	if errors.As(err, &res) && slices.Contains([]dahuarpc.ErrorType{
		dahuarpc.ErrorTypeInvalidRequest,
		dahuarpc.ErrorTypeMethodNotFound,
		dahuarpc.ErrorTypeInterfaceNotFound,
		dahuarpc.ErrorTypeUnknown,
	}, res.Type) {
		log.Err(err).Str("method", res.Method).Int("code", res.Code).Str("type", string(res.Type)).Msg("Ignoring RPC ResponseError")
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

func Normalize(ctx context.Context, db sqlite.DB) error {
	_, err := db.ExecContext(ctx, `
WITH RECURSIVE generate_series(value) AS (
  SELECT 1
  UNION ALL
  SELECT value+1 FROM generate_series WHERE value+1<=999
)
INSERT OR IGNORE INTO dahua_seeds (seed) SELECT value from generate_series;
INSERT OR IGNORE INTO dahua_event_rules (code) VALUES ('');
	`)
	if err != nil {
		return err
	}

	{
		c := NewFileCursor()
		err := db.C().DahuaNormalizeFileCursors(context.Background(), repo.DahuaNormalizeFileCursorsParams{
			QuickCursor: c.QuickCursor,
			FullCursor:  c.FullCursor,
			FullEpoch:   c.FullEpoch,
			Scan:        c.Scan,
			ScanPercent: c.ScanPercent,
			ScanType:    c.ScanType,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDahuaDetail(ctx context.Context, rpcClient dahuarpc.Conn) (models.DahuaDetail, error) {
	sn, err := magicbox.GetSerialNo(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return models.DahuaDetail{}, err
	}

	deviceClass, err := magicbox.GetDeviceClass(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return models.DahuaDetail{}, err
	}

	deviceType, err := magicbox.GetDeviceType(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return models.DahuaDetail{}, err
	}

	hardwareVersion, err := magicbox.GetHardwareVersion(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return models.DahuaDetail{}, err
	}

	marketArea, err := magicbox.GetMarketArea(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return models.DahuaDetail{}, err
	}

	processInfo, err := magicbox.GetProcessInfo(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return models.DahuaDetail{}, err
	}

	vendor, err := magicbox.GetVendor(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return models.DahuaDetail{}, err
	}

	onvifVersion, err := intervideo.ManagerGetVersion(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return models.DahuaDetail{}, err
	}

	var algorithmVersion string
	{
		res, err := peripheralchip.GetVersion(ctx, rpcClient, peripheralchip.TypeBLOB)
		if err != nil && isFatalError(err) {
			return models.DahuaDetail{}, err
		}
		if len(res) > 0 {
			algorithmVersion = res[0].SoftwareVersion
		}
	}

	return models.DahuaDetail{
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

func GetSoftwareVersion(ctx context.Context, rpcClient dahuarpc.Conn) (models.DahuaSoftwareVersion, error) {
	res, err := magicbox.GetSoftwareVersion(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return models.DahuaSoftwareVersion{}, err
	}

	return models.DahuaSoftwareVersion{
		Build:                   res.Build,
		BuildDate:               res.BuildDate,
		SecurityBaseLineVersion: res.SecurityBaseLineVersion,
		Version:                 res.Version,
		WebVersion:              res.WebVersion,
	}, nil
}

func GetLicenseList(ctx context.Context, rpcClient dahuarpc.Conn) ([]models.DahuaLicense, error) {
	licenses, err := license.GetLicenseInfo(ctx, rpcClient)
	if err != nil && isFatalError(err) {
		return nil, err
	}

	res := make([]models.DahuaLicense, 0, len(licenses))
	for _, l := range licenses {
		effectiveTime := time.Unix(int64(l.EffectiveTime), 0)

		res = append(res, models.DahuaLicense{
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

func GetStorage(ctx context.Context, rpcClient dahuarpc.Conn) ([]models.DahuaStorage, error) {
	devices, err := storage.GetDeviceAllInfo(ctx, rpcClient)
	if err != nil {
		return []models.DahuaStorage{}, checkFatalError(err)
	}

	var res []models.DahuaStorage
	for _, device := range devices {
		for _, detail := range device.Detail {
			res = append(res, models.DahuaStorage{
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

func GetError(ctx context.Context, conn dahuarpc.Client) models.DahuaError {
	err := conn.State(ctx).Error
	if err == nil {
		return models.DahuaError{}
	}

	return models.DahuaError{
		Error: err.Error(),
	}
}

func GetCoaxialStatus(ctx context.Context, rpcClient dahuarpc.Conn, channel int) (models.DahuaCoaxialStatus, error) {
	status, err := coaxialcontrolio.GetStatus(ctx, rpcClient, channel)
	if err != nil && isFatalError(err) {
		return models.DahuaCoaxialStatus{}, err
	}

	return models.DahuaCoaxialStatus{
		Speaker:    status.Speaker == "On",
		WhiteLight: status.WhiteLight == "On",
	}, nil
}

func GetCoaxialCaps(ctx context.Context, rpcClient dahuarpc.Conn, channel int) (models.DahuaCoaxialCaps, error) {
	caps, err := coaxialcontrolio.GetCaps(ctx, rpcClient, channel)
	if err != nil && isFatalError(err) {
		return models.DahuaCoaxialCaps{}, err
	}

	return models.DahuaCoaxialCaps{
		SupportControlFullcolorLight: caps.SupportControlFullcolorLight == 1,
		SupportControlLight:          caps.SupportControlLight == 1,
		SupportControlSpeaker:        caps.SupportControlSpeaker == 1,
	}, nil
}

func GetUsers(ctx context.Context, rpcClient dahuarpc.Conn, location *time.Location) ([]models.DahuaUser, error) {
	users, err := usermanager.GetActiveUserInfoAll(ctx, rpcClient)
	if err != nil {
		return nil, err
	}

	res := make([]models.DahuaUser, 0, len(users))
	for _, u := range users {
		loginTime, err := u.LoginTime.Parse(location)
		if err != nil {
			return nil, err
		}

		res = append(res, models.DahuaUser{
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

func NewDahuaEvent(v repo.DahuaEvent) models.DahuaEvent {
	return models.DahuaEvent{
		ID:        v.ID,
		DeviceID:  v.DeviceID,
		Code:      v.Code,
		Action:    v.Action,
		Index:     v.Index,
		Data:      v.Data.RawMessage,
		CreatedAt: v.CreatedAt.Time,
	}
}

func NewDahuaFile(file mediafilefind.FindNextFileInfo, affixSeed int, location *time.Location) (models.DahuaFile, error) {
	startTime, endTime, err := file.UniqueTime(affixSeed, location)
	if err != nil {
		return models.DahuaFile{}, err
	}

	return models.DahuaFile{
		Channel:     file.Channel,
		StartTime:   startTime,
		EndTime:     endTime,
		Length:      file.Length,
		Type:        file.Type,
		FilePath:    file.FilePath,
		Duration:    file.Duration,
		Disk:        file.Disk,
		VideoStream: file.VideoStream,
		Flags:       file.Flags,
		Events:      file.Events,
		Cluster:     file.Cluster,
		Partition:   file.Partition,
		PicIndex:    file.PicIndex,
		Repeat:      file.Repeat,
		WorkDir:     file.WorkDir,
		WorkDirSN:   file.WorkDirSN == 1,
		Storage:     StorageFromFilePath(file.FilePath),
	}, nil
}

func NewDahuaFiles(files []mediafilefind.FindNextFileInfo, affixSeed int, location *time.Location) ([]models.DahuaFile, error) {
	res := make([]models.DahuaFile, 0, len(files))
	for _, file := range files {
		r, err := NewDahuaFile(file, affixSeed, location)
		if err != nil {
			return []models.DahuaFile{}, err
		}

		res = append(res, r)
	}

	return res, nil
}

func GetRPCStatus(ctx context.Context, rpcClient dahuarpc.Client) models.DahuaRPCStatus {
	rpcState := rpcClient.State(ctx)
	var rpcError string
	if rpcState.Error != nil {
		rpcError = rpcState.Error.Error()
	}
	return models.DahuaRPCStatus{
		Error:     rpcError,
		State:     rpcState.State.String(),
		LastLogin: rpcState.LastLogin,
	}
}

func ListPresets(ctx context.Context, clientPTZ ptz.Client, channel int) ([]models.DahuaPreset, error) {
	vv, err := ptz.GetPresets(ctx, clientPTZ, channel)
	if err != nil {
		return nil, err
	}
	res := make([]models.DahuaPreset, 0, len(vv))
	for _, v := range vv {
		res = append(res, models.DahuaPreset{
			Index: v.Index,
			Name:  v.Name,
		})
	}
	return res, nil
}

func SetPreset(ctx context.Context, clientPTZ ptz.Client, channel, index int) error {
	return ptz.Start(ctx, clientPTZ, channel, ptz.Params{
		Code: "GotoPreset",
		Arg1: index,
	})
}

func GetUptime(ctx context.Context, c dahuarpc.Conn) (models.DahuaUptime, error) {
	uptime, err := magicbox.GetUpTime(ctx, c)
	if err != nil {
		return models.DahuaUptime{}, checkFatalError(err)
	}

	now := time.Now()

	return models.DahuaUptime{
		Last:      now.Add(-time.Duration(uptime.Last) * time.Second),
		Total:     now.Add(-time.Duration(uptime.Total) * time.Second),
		Supported: true,
	}, nil
}

func GetSunriseSunset(ctx context.Context, c dahuarpc.Conn) (models.DahuaSunriseSunset, error) {
	cfg, err := config.GetVideoInMode(ctx, c)
	if err != nil {
		return models.DahuaSunriseSunset{}, err
	}

	return models.DahuaSunriseSunset{
		SwitchMode:  cfg.Tables[0].Data.SwitchMode(),
		TimeSection: cfg.Tables[0].Data.TimeSection[0][0],
	}, nil
}

func SyncSunriseSunset(ctx context.Context, c dahuarpc.Conn, loc *time.Location, coordinate models.Coordinate, sunriseOffset, sunsetOffset time.Duration) (models.DahuaSunriseSunset, error) {
	cfg, err := config.GetVideoInMode(ctx, c)
	if err != nil {
		return models.DahuaSunriseSunset{}, err
	}

	var changed bool

	// Sync SwitchMode
	if cfg.Tables[0].Data.SwitchMode() != config.SwitchModeTime {
		cfg.Tables[0].Data.SetSwitchMode(config.SwitchModeTime)
		changed = true
	}

	// Sync TimeSection
	now := time.Now()
	sunrise, sunset := sunrise.SunriseSunset(coordinate.Latitude, coordinate.Longitude, now.Year(), now.Month(), now.Day())
	sunrise = sunrise.In(loc).Add(sunriseOffset)
	sunset = sunset.In(loc).Add(sunsetOffset)
	ts := dahuarpc.NewTimeSectionFromRange(1, sunrise, sunset)
	if cfg.Tables[0].Data.TimeSection[0][0].String() != ts.String() {
		cfg.Tables[0].Data.TimeSection[0][0] = ts
		changed = true
	}

	if changed {
		err := configmanager.SetConfig(ctx, c, cfg)
		if err != nil {
			return models.DahuaSunriseSunset{}, err
		}
	}

	return models.DahuaSunriseSunset{
		SwitchMode:  cfg.Tables[0].Data.SwitchMode(),
		TimeSection: cfg.Tables[0].Data.TimeSection[0][0],
	}, nil
}
