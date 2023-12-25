package dahuacore

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/thejerf/suture/v4"
)

type ConnRepo interface {
	GetConn(ctx context.Context, id int64) (models.DahuaConn, bool, error)
}

func NewCoaxialWorker(bus *core.Bus, deviceID int64, store *Store, repo ConnRepo) CoaxialWorker {
	return CoaxialWorker{
		bus:      bus,
		deviceID: deviceID,
		repo:     repo,
		store:    store,
	}
}

// CoaxialWorker publishes coaxial status to the bus.
type CoaxialWorker struct {
	bus      *core.Bus
	deviceID int64
	store    *Store
	repo     ConnRepo
}

func (w CoaxialWorker) String() string {
	return fmt.Sprintf("dahuacore.CoaxialWorker(id=%d)", w.deviceID)
}

func (w CoaxialWorker) Serve(ctx context.Context) error {
	conn, ok, err := w.repo.GetConn(ctx, w.deviceID)
	if err != nil {
		return err
	}
	if !ok {
		return suture.ErrDoNotRestart
	}
	client := w.store.Conn(ctx, conn)

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
