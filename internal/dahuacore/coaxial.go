package dahuacore

import (
	"context"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/thejerf/suture/v4"
)

func NewCoaxialWorker(bus *core.Bus, cameraID int64, rpcConn dahuarpc.Conn) CoaxialWorker {
	return CoaxialWorker{
		bus:      bus,
		cameraID: cameraID,
		rpcConn:  rpcConn,
	}
}

type CoaxialWorker struct {
	bus      *core.Bus
	cameraID int64
	rpcConn  dahuarpc.Conn
}

func (w CoaxialWorker) String() string {
	return fmt.Sprintf("dahuacore.CoaxialWorker(id=%d)", w.cameraID)
}

func (w CoaxialWorker) Serve(ctx context.Context) error {
	t := time.NewTicker(1 * time.Second)

	channel := 1

	// Does this camera support coaxial?
	caps, err := GetCoaxialCaps(ctx, w.cameraID, w.rpcConn, channel)
	if err != nil {
		return err
	}
	if !(caps.SupportControlSpeaker || caps.SupportControlLight || caps.SupportControlFullcolorLight) {
		return suture.ErrDoNotRestart
	}

	// Get and send initial coaxial status
	coaxialStatus, err := GetCoaxialStatus(ctx, w.cameraID, w.rpcConn, channel)
	if err != nil {
		return err
	}
	w.bus.EventDahuaCoaxialStatus(models.EventDahuaCoaxialStatus{
		Channel:       channel,
		CoaxialStatus: coaxialStatus,
	})

	// On an interval, get and send coaxial status if it changes
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}

		s, err := GetCoaxialStatus(ctx, w.cameraID, w.rpcConn, channel)
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
