package dahua

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/auth"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/coaxialcontrolio"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/intervideo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/license"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/magicbox"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/mediafilefind"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/peripheralchip"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/ptz"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/storage"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/usermanager"
	"github.com/rs/zerolog/log"
)

func NewConn(camera models.DahuaConn) Conn {
	address := NewHTTPAddress(camera.Address)
	rpcHTTPClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	cgiHTTPClient := http.Client{}

	connRPC := auth.NewClient(rpcHTTPClient, address, camera.Username, camera.Password)
	connPTZ := ptz.NewClient(connRPC)
	connCGI := dahuacgi.NewClient(cgiHTTPClient, address, camera.Username, camera.Password)

	return Conn{
		Camera: camera,
		RPC:    connRPC,
		PTZ:    connPTZ,
		CGI:    connCGI,
	}
}

type Conn struct {
	Camera models.DahuaConn
	RPC    auth.Client
	PTZ    *ptz.Client
	CGI    dahuacgi.Client
}

func ignorableError(err error) bool {
	res := &dahuarpc.ResponseError{}
	if errors.As(err, &res) && slices.Contains([]dahuarpc.ResponseErrorType{
		dahuarpc.ErrResponseTypeInvalidRequest,
		dahuarpc.ErrResponseTypeMethodNotFound,
		dahuarpc.ErrResponseTypeInterfaceNotFound,
		dahuarpc.ErrResponseTypeUnknown,
	}, res.Type) {
		log.Err(err).Str("method", res.Method).Int("code", res.Code).Str("type", string(res.Type)).Msg("Ignoring ResponseError")
		return true
	}

	return false
}

func GetDahuaDetail(ctx context.Context, cameraID int64, rpcClient dahuarpc.Conn) (models.DahuaDetail, error) {
	sn, err := magicbox.GetSerialNo(ctx, rpcClient)
	if err != nil && !ignorableError(err) {
		return models.DahuaDetail{}, err
	}

	deviceClass, err := magicbox.GetDeviceClass(ctx, rpcClient)
	if err != nil && !ignorableError(err) {
		return models.DahuaDetail{}, err
	}

	deviceType, err := magicbox.GetDeviceType(ctx, rpcClient)
	if err != nil && !ignorableError(err) {
		return models.DahuaDetail{}, err
	}

	hardwareVersion, err := magicbox.GetHardwareVersion(ctx, rpcClient)
	if err != nil && !ignorableError(err) {
		return models.DahuaDetail{}, err
	}

	marketArea, err := magicbox.GetMarketArea(ctx, rpcClient)
	if err != nil && !ignorableError(err) {
		return models.DahuaDetail{}, err
	}

	ProcessInfo, err := magicbox.GetProcessInfo(ctx, rpcClient)
	if err != nil && !ignorableError(err) {
		return models.DahuaDetail{}, err
	}

	vendor, err := magicbox.GetVendor(ctx, rpcClient)
	if err != nil && !ignorableError(err) {
		return models.DahuaDetail{}, err
	}

	onvifVersion, err := intervideo.ManagerGetVersion(ctx, rpcClient)
	if err != nil && !ignorableError(err) {
		return models.DahuaDetail{}, err
	}

	var algorithmVersion string
	{
		res, err := peripheralchip.GetVersion(ctx, rpcClient, peripheralchip.TypeBLOB)
		if err != nil && !ignorableError(err) {
			return models.DahuaDetail{}, err
		}
		if len(res) > 0 {
			algorithmVersion = res[0].SoftwareVersion
		}
	}

	return models.DahuaDetail{
		CameraID:         cameraID,
		SN:               sn,
		DeviceClass:      deviceClass,
		DeviceType:       deviceType,
		HardwareVersion:  hardwareVersion,
		MarketArea:       marketArea,
		ProcessInfo:      ProcessInfo,
		Vendor:           vendor,
		OnvifVersion:     onvifVersion,
		AlgorithmVersion: algorithmVersion,
	}, nil
}

func GetSoftwareVersion(ctx context.Context, cameraID int64, rpcClient dahuarpc.Conn) (models.DahuaSoftwareVersion, error) {
	res, err := magicbox.GetSoftwareVersion(ctx, rpcClient)
	if err != nil && !ignorableError(err) {
		return models.DahuaSoftwareVersion{}, err
	}

	return models.DahuaSoftwareVersion{
		CameraID:                cameraID,
		Build:                   res.Build,
		BuildDate:               res.BuildDate,
		SecurityBaseLineVersion: res.SecurityBaseLineVersion,
		Version:                 res.Version,
		WebVersion:              res.WebVersion,
	}, nil
}

func GetLicenseList(ctx context.Context, cameraID int64, rpcClient dahuarpc.Conn) ([]models.DahuaLicense, error) {
	licenses, err := license.GetLicenseInfo(ctx, rpcClient)
	if err != nil && !ignorableError(err) {
		return nil, err
	}

	res := make([]models.DahuaLicense, 0, len(licenses))
	for _, l := range licenses {
		effectiveTime := time.Unix(int64(l.EffectiveTime), 0)

		res = append(res, models.DahuaLicense{
			CameraID:      cameraID,
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

func GetStorage(ctx context.Context, cameraID int64, rpcClient dahuarpc.Conn) ([]models.DahuaStorage, error) {
	devices, err := storage.GetDeviceAllInfo(ctx, rpcClient)
	if err != nil {
		if ignorableError(err) {
			return []models.DahuaStorage{}, nil
		}
		return nil, err
	}

	var res []models.DahuaStorage
	for _, device := range devices {
		for _, detail := range device.Detail {
			res = append(res, models.DahuaStorage{
				CameraID:   cameraID,
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

func GetError(conn auth.Client) models.Error {
	err := conn.Data().Error
	if err == nil {
		return models.Error{}
	}

	return models.Error{
		Error: err.Error(),
	}
}

func GetCoaxialStatus(ctx context.Context, cameraID int64, rpcClient dahuarpc.Conn, channel int) (models.DahuaCoaxialStatus, error) {
	if channel == 0 {
		channel = 1
	}

	status, err := coaxialcontrolio.GetStatus(ctx, rpcClient, channel)
	if err != nil {
		return models.DahuaCoaxialStatus{}, err
	}

	return models.DahuaCoaxialStatus{
		CameraID:   cameraID,
		Speaker:    status.Speaker == "On",
		WhiteLight: status.WhiteLight == "On",
	}, nil
}

func GetCoaxialCaps(ctx context.Context, cameraID int64, rpcClient dahuarpc.Conn, channel int) (models.DahuaCoaxialCaps, error) {
	if channel == 0 {
		channel = 1
	}

	caps, err := coaxialcontrolio.GetCaps(ctx, rpcClient, channel)
	if err != nil {
		return models.DahuaCoaxialCaps{}, err
	}

	return models.DahuaCoaxialCaps{
		CameraID:                     cameraID,
		SupportControlFullcolorLight: caps.SupportControlFullcolorLight == 1,
		SupportControlLight:          caps.SupportControlLight == 1,
		SupportControlSpeaker:        caps.SupportControlSpeaker == 1,
	}, nil
}

func GetUsers(ctx context.Context, cameraID int64, rpcClient dahuarpc.Conn, location *time.Location) ([]models.DahuaUser, error) {
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
			CameraID:      cameraID,
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

func NewDahuaEvent(cameraID int64, event dahuacgi.Event, createdAt time.Time) models.DahuaEvent {
	return models.DahuaEvent{
		CameraID:  cameraID,
		Code:      event.Code,
		Action:    event.Action,
		Index:     event.Index,
		Data:      event.Data,
		CreatedAt: createdAt,
	}
}

func NewDahuaFile(cameraID int64, file mediafilefind.FindNextFileInfo, affixSeed int, location *time.Location) (models.DahuaFile, error) {
	startTime, endTime, err := file.UniqueTime(affixSeed, location)
	if err != nil {
		return models.DahuaFile{}, err
	}

	return models.DahuaFile{
		CameraID:    cameraID,
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
		WorkDirSN:   file.WorkDirSN,
	}, nil
}

func NewDahuaFiles(cameraID int64, files []mediafilefind.FindNextFileInfo, affixSeed int, location *time.Location) ([]models.DahuaFile, error) {
	res := make([]models.DahuaFile, 0, len(files))
	for _, file := range files {
		r, err := NewDahuaFile(cameraID, file, affixSeed, location)
		if err != nil {
			return []models.DahuaFile{}, err
		}

		res = append(res, r)
	}

	return res, nil
}

func NewDahuaCamera(c models.DahuaCamera) (models.DahuaCamera, error) {
	if c.Location.Location == nil {
		c.Location = types.Location{Location: time.Local}
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	return c, validate.Validate.Struct(c)
}

func UpdateDahuaCamera(c models.DahuaCamera) (models.DahuaCamera, error) {
	c.UpdatedAt = time.Now()
	return c, validate.Validate.Struct(c)
}

func GetSeed(c models.DahuaConn) int {
	if c.Seed != 0 {
		return c.Seed
	}

	return int(c.ID)
}

func NewHTTPAddress(address string) string {
	return "http://" + address
}

func GetDahuaStatus(camera models.DahuaConn, rpcClient auth.Client) models.DahuaStatus {
	rpcData := rpcClient.Data()
	var rpcError string
	if rpcData.Error != nil {
		rpcError = rpcData.Error.Error()
	}
	return models.DahuaStatus{
		CameraID:     camera.ID,
		Address:      camera.Address,
		Username:     camera.Username,
		Location:     camera.Location.String(),
		Seed:         camera.Seed,
		RPCError:     rpcError,
		RPCState:     rpcData.State.String(),
		RPCLastLogin: rpcData.LastLogin,
	}
}

func SetPreset(ctx context.Context, clientPTZ *ptz.Client, channel, index int) error {
	return ptz.Start(ctx, clientPTZ, channel, ptz.Params{
		Code: "GotoPreset",
		Arg1: index,
	})
}

type CameraStore interface {
	List(ctx context.Context) ([]models.DahuaConn, error)
}

func Bootstrap(ctx context.Context, cameraStore CameraStore, store *Store, eventWorkerStore *EventWorkerStore) error {
	cameras, err := cameraStore.List(ctx)
	if err != nil {
		return err
	}
	conns := store.ConnList(ctx, cameras)
	for _, conn := range conns {
		if err := eventWorkerStore.Create(conn.Camera); err != nil {
			return err
		}
	}
	return err
}
