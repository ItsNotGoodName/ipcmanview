package dahua

import (
	"context"
	"net"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/internal/system/action"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

func deviceFrom(v repo.DahuaDevice) _Device {
	return _Device{
		Name:     v.Name,
		URL:      v.Url.URL,
		Username: v.Username,
		Email:    v.Email.String,
	}
}

type _Device struct {
	Name     string `validate:"required,lte=64"`
	URL      *url.URL
	Username string
	Email    string `validate:"omitempty,lte=128,email"`
}

func (d *_Device) normalize(create bool) {
	// Name
	d.Name = strings.TrimSpace(d.Name)
	// Email
	d.Email = strings.TrimSpace(d.Email)
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
			d.Name = d.URL.Hostname()
		}
		if d.Username == "" {
			d.Username = "admin"
		}
	}
}

func (d *_Device) getIP() (string, error) {
	ip := d.URL.Hostname()

	ips, err := net.LookupIP(ip)
	if err != nil {
		return "", core.NewFieldError("URL", err.Error())
	}

	for _, v := range ips {
		if v.To4() != nil {
			ip = v.String()
			break
		}
	}

	return ip, nil
}

type CreateDeviceParams struct {
	Name     string
	URL      *url.URL
	Username string
	Password string
	Location *time.Location
	Feature  models.DahuaFeature
	Email    string
}

func CreateDevice(ctx context.Context, arg CreateDeviceParams) (int64, error) {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return 0, err
	}

	model := _Device{
		Name:     arg.Name,
		URL:      arg.URL,
		Username: arg.Username,
		Email:    arg.Email,
	}
	model.normalize(true)

	err := core.ValidateStruct(ctx, model)
	if err != nil {
		return 0, err
	}

	ip, err := model.getIP()
	if err != nil {
		return 0, err
	}

	now := types.NewTime(time.Now())
	id, err := createDahuaDevice(ctx, repo.DahuaCreateDeviceParams{
		Name:      model.Name,
		Url:       types.NewURL(model.URL),
		Ip:        ip,
		Username:  model.Username,
		Password:  arg.Password,
		Location:  types.NewLocation(arg.Location),
		Feature:   arg.Feature,
		Email:     core.StringToNullString(model.Email),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return 0, err
	}

	return id, err
}

func createDahuaDevice(ctx context.Context, arg repo.DahuaCreateDeviceParams) (int64, error) {
	tx, err := app.DB.BeginTx(ctx, true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	id, err := tx.C().DahuaCreateDevice(ctx, arg)
	if err != nil {
		return 0, err
	}

	if err := tx.C().DahuaAllocateSeed(ctx, core.NewNullInt64(id)); err != nil {
		return 0, err
	}

	arg2 := newFileCursor()
	arg2.DeviceID = id
	if err := tx.C().DahuaCreateFileCursor(ctx, arg2); err != nil {
		return 0, err
	}

	if err := system.CreateEvent(ctx, tx.C(), action.DahuaDeviceCreated.Create(id)); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	app.Hub.DahuaDeviceCreated(bus.DahuaDeviceCreated{
		DeviceID: id,
	})

	return id, nil
}

type UpdateDeviceParams struct {
	ID          int64
	Name        string
	URL         *url.URL
	Username    string
	NewPassword string
	Location    *time.Location
	Feature     models.DahuaFeature
	Email       string
}

func UpdateDevice(ctx context.Context, arg UpdateDeviceParams) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	dbModel, err := GetDevice(ctx, arg.ID)
	if err != nil {
		return err
	}
	model := deviceFrom(dbModel)

	// Mutate
	model.Name = arg.Name
	model.URL = arg.URL
	model.Username = arg.Username
	model.Email = arg.Email
	model.normalize(false)

	if err := core.ValidateStruct(ctx, model); err != nil {
		return err
	}

	ip, err := model.getIP()
	if err != nil {
		return err
	}

	password := dbModel.Password
	if arg.NewPassword != "" {
		password = arg.NewPassword
	}

	return updateDevice(ctx, repo.DahuaUpdateDeviceParams{
		Name:      model.Name,
		Url:       types.NewURL(model.URL),
		Ip:        ip,
		Username:  model.Username,
		Password:  password,
		Location:  types.NewLocation(arg.Location),
		Feature:   arg.Feature,
		Email:     core.StringToNullString(model.Email),
		UpdatedAt: types.NewTime(time.Now()),
		ID:        dbModel.ID,
	})
}

func updateDevice(ctx context.Context, arg repo.DahuaUpdateDeviceParams) error {
	tx, err := app.DB.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.C().DahuaUpdateDevice(ctx, arg); err != nil {
		return err
	}

	if err := system.CreateEvent(ctx, tx.C(), action.DahuaDeviceUpdated.Create(arg.ID)); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	app.Hub.DahuaDeviceUpdated(bus.DahuaDeviceUpdated{
		DeviceID: arg.ID,
	})

	return nil
}

func DeleteDevice(ctx context.Context, id int64) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	return deleteDevice(ctx, id)
}

func deleteDevice(ctx context.Context, id int64) error {
	tx, err := app.DB.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.C().DahuaDeleteDevice(ctx, id); err != nil {
		return err
	}

	if err := system.CreateEvent(ctx, tx.C(), action.DahuaDeviceDeleted.Create(id)); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	app.Hub.DahuaDeviceDeleted(bus.DahuaDeviceDeleted{
		DeviceID: id,
	})

	return err
}

func UpdateDeviceDisabled(ctx context.Context, id int64, disable bool) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	return updateDeviceDisabled(ctx, repo.DahuaUpdateDeviceDisabledAtParams{
		DisabledAt: types.NullTime{
			Time:  types.NewTime(time.Now()),
			Valid: disable,
		},
		ID: id,
	})
}

func updateDeviceDisabled(ctx context.Context, arg repo.DahuaUpdateDeviceDisabledAtParams) error {
	tx, err := app.DB.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.C().DahuaUpdateDeviceDisabledAt(ctx, arg); err != nil {
		return err
	}

	if err := system.CreateEvent(ctx, tx.C(), action.DahuaDeviceUpdated.Create(arg.ID)); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	app.Hub.DahuaDeviceUpdated(bus.DahuaDeviceUpdated{
		DeviceID: arg.ID,
	})

	return nil
}
