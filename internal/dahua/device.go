package dahua

import (
	"context"
	"net"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/common"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/internal/validate"
)

func normalizeDeviceHTTPURL(u *url.URL) *url.URL {
	if slices.Contains([]string{"http", "https"}, u.Scheme) {
		return u
	}

	switch u.Port() {
	case "443":
		u.Scheme = "https"
	default:
		u.Scheme = "http"
	}

	u, err := url.Parse(u.String())
	if err != nil {
		panic(err)
	}

	return u
}

func normalizeDevice(arg models.DahuaDevice, create bool) models.DahuaDevice {
	arg.Name = strings.TrimSpace(arg.Name)
	arg.Url = normalizeDeviceHTTPURL(arg.Url)

	if create {
		if arg.Name == "" {
			arg.Name = arg.Url.String()
		}
		if arg.Username == "" {
			arg.Username = "admin"
		}
	}

	return arg
}

func parseDeviceIPFromURL(urL *url.URL) (string, error) {
	ip := urL.Hostname()

	ips, err := net.LookupIP(ip)
	if err != nil {
		return "", err
	}

	for _, i2 := range ips {
		if i2.To4() != nil {
			ip = i2.String()
			break
		}
	}

	return ip, nil
}

func CreateDevice(ctx context.Context, db repo.DB, bus *common.Bus, arg models.DahuaDevice) (models.DahuaDeviceConn, error) {
	arg = normalizeDevice(arg, true)

	err := validate.Validate.Struct(arg)
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}

	ip, err := parseDeviceIPFromURL(arg.Url)
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}

	now := types.NewTime(time.Now())
	id, err := db.CreateDahuaDevice(ctx, repo.CreateDahuaDeviceParams{
		Name:      arg.Name,
		Url:       types.NewURL(arg.Url),
		Ip:        ip,
		Username:  arg.Username,
		Password:  arg.Password,
		Location:  types.NewLocation(arg.Location),
		Feature:   arg.Feature,
		CreatedAt: now,
		UpdatedAt: now,
	}, NewFileCursor())
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}

	dbDevice, err := db.GetDahuaDevice(ctx, id)
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}
	device := dbDevice.Convert()

	bus.EventDahuaDeviceCreated(models.EventDahuaDeviceCreated{
		Device: device,
	})

	return device, err
}

func UpdateDevice(ctx context.Context, db repo.DB, bus *common.Bus, arg models.DahuaDevice) (models.DahuaDeviceConn, error) {
	arg = normalizeDevice(arg, false)

	err := validate.Validate.Struct(arg)
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}

	ip, err := parseDeviceIPFromURL(arg.Url)
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}

	_, err = db.UpdateDahuaDevice(ctx, repo.UpdateDahuaDeviceParams{
		Name:      arg.Name,
		Url:       types.NewURL(arg.Url),
		Ip:        ip,
		Username:  arg.Username,
		Password:  arg.Password,
		Location:  types.NewLocation(arg.Location),
		Feature:   arg.Feature,
		UpdatedAt: types.NewTime(time.Now()),
		ID:        arg.ID,
	})
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}

	dbDevice, err := db.GetDahuaDevice(ctx, arg.ID)
	if err != nil {
		return models.DahuaDeviceConn{}, err
	}
	device := dbDevice.Convert()

	bus.EventDahuaDeviceUpdated(models.EventDahuaDeviceUpdated{
		Device: device,
	})

	return device, nil
}

func DeleteDevice(ctx context.Context, db repo.DB, bus *common.Bus, id int64) error {
	if err := db.DeleteDahuaDevice(ctx, id); err != nil {
		return err
	}
	bus.EventDahuaDeviceDeleted(models.EventDahuaDeviceDeleted{
		DeviceID: id,
	})
	return nil
}
