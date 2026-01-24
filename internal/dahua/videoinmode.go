package dahua

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager/config"
	"github.com/jmoiron/sqlx"
	"github.com/nathan-osman/go-sunrise"
)

type DeviceVideoInMode struct {
	SwitchMode  string `json:"switch_mode"`
	TimeSection string `json:"time_section"`
}

func GetVideoInMode(ctx context.Context, c dahuarpc.Conn) (DeviceVideoInMode, error) {
	cfg, err := config.GetVideoInMode(ctx, c)
	if err != nil {
		return DeviceVideoInMode{}, err
	}

	return DeviceVideoInMode{
		SwitchMode:  cfg.Tables[0].Data.SwitchMode().String(),
		TimeSection: cfg.Tables[0].Data.TimeSection[0][0].String(),
	}, nil
}

type SyncVideoInModeArgs struct {
	Location      *time.Location
	Latitude      float64
	Longitude     float64
	SunriseOffset time.Duration
	SunsetOffset  time.Duration
}

func SyncVideoInMode(ctx context.Context, c dahuarpc.Conn, args SyncVideoInModeArgs) (DeviceVideoInMode, error) {
	cfg, err := config.GetVideoInMode(ctx, c)
	if err != nil {
		return DeviceVideoInMode{}, err
	}

	var changed bool

	// Sync SwitchMode
	if cfg.Tables[0].Data.SwitchMode() != config.SwitchModeSchedule {
		cfg.Tables[0].Data.SetSwitchMode(config.SwitchModeSchedule)
		changed = true
	}

	// Sync TimeSection
	now := time.Now()
	sunrise, sunset := sunrise.SunriseSunset(args.Latitude, args.Longitude, now.Year(), now.Month(), now.Day())
	sunrise = sunrise.In(args.Location).Add(args.SunriseOffset)
	sunset = sunset.In(args.Location).Add(args.SunsetOffset)
	ts := dahuarpc.NewTimeSectionFromRange(1, sunrise, sunset)
	if cfg.Tables[0].Data.TimeSection[0][0].String() != ts.String() {
		cfg.Tables[0].Data.TimeSection[0][0] = ts
		changed = true
	}

	if changed {
		err := configmanager.SetConfigArray(ctx, c, cfg)
		if err != nil {
			return DeviceVideoInMode{}, err
		}
	}

	return DeviceVideoInMode{
		SwitchMode:  cfg.Tables[0].Data.SwitchMode().String(),
		TimeSection: cfg.Tables[0].Data.TimeSection[0][0].String(),
	}, nil
}

func NewSyncVideoInModeJob(db *sqlx.DB, store *Store) SyncVideoInModeJob {
	return SyncVideoInModeJob{
		db:    db,
		store: store,
	}
}

type SyncVideoInModeJob struct {
	db    *sqlx.DB
	store *Store
}

func (w SyncVideoInModeJob) Description() string {
	return "dahua.SyncVideoInModeJob"
}

func (w SyncVideoInModeJob) Execute(ctx context.Context) error {
	var devices []Device
	err := w.db.Select(&devices, `
		SELECT d.* 
		FROM dahua_devices AS d, (SELECT value AS sync_video_in_mode FROM settings WHERE key = ?) AS s
		WHERE coalesce(d.sync_video_in_mode, s.sync_video_in_mode) IS TRUE
	`, system.KeySyncVideoInMode)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for _, device := range devices {
		slog := slog.With("device", device.Name, "service", w.Description())

		wg.Add(1)
		go func(device Device) {
			defer wg.Done()

			position, err := GetDevicePosition(ctx, w.db, device.ID)
			if err != nil {
				slog.Error("Failed to get device position", "error", err)
				return
			}

			client, err := w.store.GetClient(ctx, device.Key)
			if err != nil {
				slog.Error("Failed to get client", "error", err)
				return
			}

			_, err = SyncVideoInMode(ctx, client.RPC, SyncVideoInModeArgs{
				Location:      position.Location.Location,
				Latitude:      position.Latitude,
				Longitude:     position.Longitude,
				SunriseOffset: position.Sunrise_Offset.Duration,
				SunsetOffset:  position.Sunset_Offset.Duration,
			})
			if err != nil {
				slog.Error("Failed to sync video in mode", "error", err)
				return
			}
		}(device)
	}

	wg.Wait()

	return nil
}
