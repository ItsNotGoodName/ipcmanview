package dahua

import (
	"context"
	"net"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/event/action"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
)

func deviceFrom(v repo.DahuaDevice) device {
	return device{
		Name:     v.Name,
		URL:      v.Url.URL,
		Username: v.Username,
	}
}

type device struct {
	Name     string `validate:"required,lte=64"`
	URL      *url.URL
	Username string
}

func (d *device) normalize(create bool) {
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
			d.Name = d.URL.Hostname()
		}
		if d.Username == "" {
			d.Username = "admin"
		}
	}
}

func (d *device) getIP() (string, error) {
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
}

func CreateDevice(ctx context.Context, db sqlite.DB, bus *event.Bus, arg CreateDeviceParams) (int64, error) {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return 0, err
	}

	model := device{
		Name:     arg.Name,
		URL:      arg.URL,
		Username: arg.Username,
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
	id, err := createDahuaDevice(ctx, db, bus, repo.DahuaCreateDeviceParams{
		Name:      model.Name,
		Url:       types.NewURL(model.URL),
		Ip:        ip,
		Username:  model.Username,
		Password:  arg.Password,
		Location:  types.NewLocation(arg.Location),
		Feature:   arg.Feature,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return 0, err
	}

	return id, err
}

func createDahuaDevice(ctx context.Context, db sqlite.DB, bus *event.Bus, arg repo.DahuaCreateDeviceParams) (int64, error) {
	tx, err := db.BeginTx(ctx, true)
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

	if err := event.CreateEventAndCommit(ctx, tx, action.DahuaDeviceCreated.Create(id)); err != nil {
		return 0, err
	}
	bus.DahuaDeviceCreated(event.DahuaDeviceCreated{
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
}

func UpdateDevice(ctx context.Context, db sqlite.DB, bus *event.Bus, arg UpdateDeviceParams) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	dbModel, err := GetDevice(ctx, db, GetDeviceFilter{
		ID: arg.ID,
	})
	if err != nil {
		return err
	}
	model := deviceFrom(dbModel)

	// Mutate
	model.Name = arg.Name
	model.URL = arg.URL
	model.Username = arg.Username
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

	return updateDevice(ctx, db, bus, repo.DahuaUpdateDeviceParams{
		Name:      model.Name,
		Url:       types.NewURL(model.URL),
		Ip:        ip,
		Username:  model.Username,
		Password:  password,
		Location:  types.NewLocation(arg.Location),
		Feature:   arg.Feature,
		UpdatedAt: types.NewTime(time.Now()),
		ID:        dbModel.ID,
	})
}

func updateDevice(ctx context.Context, db sqlite.DB, bus *event.Bus, arg repo.DahuaUpdateDeviceParams) error {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.C().DahuaUpdateDevice(ctx, arg); err != nil {
		return err
	}

	if err := event.CreateEventAndCommit(ctx, tx, action.DahuaDeviceUpdated.Create(arg.ID)); err != nil {
		return err
	}
	bus.DahuaDeviceUpdated(event.DahuaDeviceUpdated{
		DeviceID: arg.ID,
	})

	return nil
}

func DeleteDevice(ctx context.Context, db sqlite.DB, bus *event.Bus, id int64) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	return deleteDevice(ctx, db, bus, id)
}

func deleteDevice(ctx context.Context, db sqlite.DB, bus *event.Bus, id int64) error {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.C().DahuaDeleteDevice(ctx, id); err != nil {
		return err
	}

	if err := event.CreateEventAndCommit(ctx, tx, action.DahuaDeviceDeleted.Create(id)); err != nil {
		return err
	}
	bus.DahuaDeviceDeleted(event.DahuaDeviceDeleted{
		DeviceID: id,
	})

	return err
}

func UpdateDeviceDisabled(ctx context.Context, db sqlite.DB, bus *event.Bus, id int64, disable bool) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	if disable {
		return updateDeviceDisabled(ctx, db, bus, repo.DahuaUpdateDeviceDisabledAtParams{
			DisabledAt: types.NewNullTime(time.Now()),
			ID:         id,
		})
	}
	return updateDeviceDisabled(ctx, db, bus, repo.DahuaUpdateDeviceDisabledAtParams{
		DisabledAt: types.NullTime{},
		ID:         id,
	})
}

func updateDeviceDisabled(ctx context.Context, db sqlite.DB, bus *event.Bus, arg repo.DahuaUpdateDeviceDisabledAtParams) error {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.C().DahuaUpdateDeviceDisabledAt(ctx, arg); err != nil {
		return err
	}

	if err := event.CreateEventAndCommit(ctx, tx, action.DahuaDeviceUpdated.Create(arg.ID)); err != nil {
		return err
	}
	bus.DahuaDeviceUpdated(event.DahuaDeviceUpdated{
		DeviceID: arg.ID,
	})

	return nil
}
