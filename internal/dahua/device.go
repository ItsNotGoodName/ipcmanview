package dahua

import (
	"context"
	"database/sql"
	"net"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

func NewDevice(v repo.DahuaDevice) Device {
	return Device{
		ID:       v.ID,
		Name:     v.Name,
		URL:      v.Url.URL,
		Username: v.Username,
		Password: v.Password,
		Location: v.Location.Location,
		Feature:  v.Feature,
	}
}

type Device struct {
	ID          int64
	Name        string `validate:"required,lte=64"`
	URL         *url.URL
	Username    string
	Password    string
	NewPassword string
	Location    *time.Location
	Feature     models.DahuaFeature
}

func (d *Device) normalize(create bool) {
	// Name
	d.Name = strings.TrimSpace(d.Name)
	// URL
	if !slices.Contains([]string{"http", "https"}, d.URL.Scheme) {
		switch d.URL.Port() {
		case "443":
			d.URL.Scheme = "https"
		default:
			d.URL.Scheme = "http"
		}

		u, err := url.Parse(d.URL.String())
		if err != nil {
			panic(err)
		}
		d.URL = u
	}

	// Name/Username
	if create {
		if d.Name == "" {
			d.Name = d.URL.String()
		}
		if d.Username == "" {
			d.Username = "admin"
		}
	}
}

func (d *Device) getIP() (string, error) {
	ip := d.URL.Hostname()

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

func createDahuaDevice(ctx context.Context, db repo.DB, arg repo.DahuaCreateDeviceParams) (int64, error) {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	id, err := tx.DahuaCreateDevice(ctx, arg)
	if err != nil {
		return 0, err
	}

	// TODO: sql.NullInt64 should just be int64
	err = tx.DahuaAllocateSeed(ctx, sql.NullInt64{
		Valid: true,
		Int64: id,
	})
	if err != nil {
		return 0, err
	}

	arg2 := NewFileCursor()
	arg2.DeviceID = id
	err = tx.DahuaCreateFileCursor(ctx, arg2)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func CreateDevice(ctx context.Context, db repo.DB, bus *event.Bus, arg Device) (int64, error) {
	arg.normalize(true)

	err := core.Validate.Struct(arg)
	if err != nil {
		return 0, err
	}

	ip, err := arg.getIP()
	if err != nil {
		return 0, err
	}

	now := types.NewTime(time.Now())
	id, err := createDahuaDevice(ctx, db, repo.DahuaCreateDeviceParams{
		Name:      arg.Name,
		Url:       types.NewURL(arg.URL),
		Ip:        ip,
		Username:  arg.Username,
		Password:  arg.Password,
		Location:  types.NewLocation(arg.Location),
		Feature:   arg.Feature,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return 0, err
	}

	device, err := db.DahuaGetDevice(ctx, repo.FatDahuaDeviceParams{IDs: []int64{id}})
	if err != nil {
		return 0, err
	}

	bus.DahuaDeviceCreated(event.DahuaDeviceCreated{
		Device: device,
	})

	return id, err
}

func UpdateDevice(ctx context.Context, db repo.DB, bus *event.Bus, arg Device) error {
	arg.normalize(true)

	err := core.Validate.Struct(arg)
	if err != nil {
		return err
	}

	ip, err := arg.getIP()
	if err != nil {
		return err
	}

	password := arg.Password
	if arg.NewPassword != "" {
		password = arg.NewPassword
	}

	_, err = db.DahuaUpdateDevice(ctx, repo.DahuaUpdateDeviceParams{
		Name:      arg.Name,
		Url:       types.NewURL(arg.URL),
		Ip:        ip,
		Username:  arg.Username,
		Password:  password,
		Location:  types.NewLocation(arg.Location),
		Feature:   arg.Feature,
		UpdatedAt: types.NewTime(time.Now()),
		ID:        arg.ID,
	})
	if err != nil {
		return err
	}

	device, err := db.DahuaGetDevice(ctx, repo.FatDahuaDeviceParams{IDs: []int64{arg.ID}})
	if err != nil {
		return err
	}

	bus.DahuaDeviceUpdated(event.DahuaDeviceUpdated{
		Device: device,
	})

	return nil
}

func DeleteDevice(ctx context.Context, db repo.DB, bus *event.Bus, id int64) error {
	if err := db.DahuaDeleteDevice(ctx, id); err != nil {
		return err
	}
	bus.DahuaDeviceDeleted(event.DahuaDeviceDeleted{
		DeviceID: id,
	})
	return nil
}
