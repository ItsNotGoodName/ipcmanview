package dahua

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/thejerf/suture/v4"
)

func NewCoaxialWorker(bus *core.Bus, deviceID int64, store *Store, db repo.DB) CoaxialWorker {
	return CoaxialWorker{
		bus:      bus,
		deviceID: deviceID,
		db:       db,
		store:    store,
	}
}

// CoaxialWorker publishes coaxial status to the bus.
type CoaxialWorker struct {
	bus      *core.Bus
	deviceID int64
	store    *Store
	db       repo.DB
}

func (w CoaxialWorker) String() string {
	return fmt.Sprintf("dahua.CoaxialWorker(id=%d)", w.deviceID)
}

func (w CoaxialWorker) Serve(ctx context.Context) error {
	return sutureext.SanitizeError(ctx, w.serve(ctx))
}

func (w CoaxialWorker) serve(ctx context.Context) error {
	dbDevice, err := w.db.GetDahuaDevice(ctx, w.deviceID)
	if err != nil {
		if repo.IsNotFound(err) {
			return suture.ErrDoNotRestart
		}
		return err
	}
	client := w.store.Client(ctx, dbDevice.Convert().DahuaConn)

	channel := 1

	// Does this device support coaxial?
	caps, err := GetCoaxialCaps(ctx, w.deviceID, client.RPC, channel)
	if err != nil {
		return err
	}
	if !(caps.SupportControlSpeaker || caps.SupportControlLight || caps.SupportControlFullcolorLight) {
		return suture.ErrDoNotRestart
	}

	// Get and send initial coaxial status
	coaxialStatus, err := GetCoaxialStatus(ctx, w.deviceID, client.RPC, channel)
	if err != nil {
		return err
	}
	w.bus.EventDahuaCoaxialStatus(models.EventDahuaCoaxialStatus{
		Channel:       channel,
		CoaxialStatus: coaxialStatus,
	})

	t := time.NewTicker(1 * time.Second)

	// Get and send coaxial status if it changes on an interval
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}

		s, err := GetCoaxialStatus(ctx, w.deviceID, client.RPC, channel)
		if err != nil {
			return err
		}
		if coaxialStatus.Speaker == s.Speaker && coaxialStatus.WhiteLight == s.WhiteLight {
			continue
		}
		coaxialStatus = s

		w.bus.EventDahuaCoaxialStatus(models.EventDahuaCoaxialStatus{
			Channel:       channel,
			CoaxialStatus: coaxialStatus,
		})
	}
}
