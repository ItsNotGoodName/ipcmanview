package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/encode"
)

func UpdateStream(ctx context.Context, db repo.DB, stream repo.DahuaStream, arg repo.UpdateDahuaStreamParams) (repo.DahuaStream, error) {
	return db.UpdateDahuaStream(ctx, arg)
}

func SupportStreams(device models.DahuaDeviceConn) bool {
	return device.DahuaConn.Feature.EQ(models.DahuaFeatureCamera)
}

// SyncStreams fetches streams from device and inserts them into the database.
// SupportStreams should be called to check if sync streams is possible.
func SyncStreams(ctx context.Context, db repo.DB, deviceID int64, conn dahuarpc.Conn) error {
	caps, err := encode.GetCaps(ctx, conn, 1)
	if err != nil {
		return err
	}

	subtypes := 1
	if caps.MaxExtraStream > 0 && caps.MaxExtraStream < 10 {
		subtypes += caps.MaxExtraStream
	}

	for channelIndex, device := range caps.VideoEncodeDevices {
		names := make([]string, subtypes)
		for i, v := range device.SupportDynamicBitrate {
			if i < len(names) {
				names[i] = v.Stream
			}
		}

		args := repo.TryCreateDahuaStreamParams{
			DeviceID: deviceID,
			Channel:  int64(channelIndex + 1),
		}
		for i := 0; i < subtypes; i++ {
			args.Subtype = int64(i)
			args.Name = names[i]
			err := db.TryCreateDahuaStream(ctx, args)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func RegisterStreams(bus *core.Bus, db repo.DB, store *Store) {
	bus.OnEventDahuaDeviceCreated(func(ctx context.Context, event models.EventDahuaDeviceCreated) error {
		if SupportStreams(event.Device) {
			return SyncStreams(ctx, db, event.Device.DahuaConn.ID, store.Client(ctx, event.Device.DahuaConn).RPC)
		}

		return nil
	})
	bus.OnEventDahuaDeviceUpdated(func(ctx context.Context, event models.EventDahuaDeviceUpdated) error {
		if SupportStreams(event.Device) {
			return SyncStreams(ctx, db, event.Device.DahuaConn.ID, store.Client(ctx, event.Device.DahuaConn).RPC)
		}

		return nil
	})
}
