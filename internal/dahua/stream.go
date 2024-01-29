package dahua

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/ipcmanview/internal/event"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/encode"
)

type createInternalStreamsParams struct {
	Channel int64
	Subtype int64
	Name    string
}

func createInternalStreams(ctx context.Context, db repo.DB, deviceID int64, args []createInternalStreamsParams) error {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.DahuaUpdateStreamForInternal(ctx, deviceID)
	if err != nil {
		return err
	}

	ids := make([]int64, 0, len(args))
	for _, arg := range args {
		id, err := tx.DahuaCreateStreamForInternal(ctx, repo.DahuaCreateStreamForInternalParams{
			DeviceID: deviceID,
			Channel:  arg.Channel,
			Subtype:  arg.Subtype,
			Name:     arg.Name,
		})
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	return tx.Commit()
}

func UpdateStream(ctx context.Context, db repo.DB, stream repo.DahuaStream, arg repo.DahuaUpdateStreamParams) (repo.DahuaStream, error) {
	return db.DahuaUpdateStream(ctx, arg)
}

func DeleteStream(ctx context.Context, db repo.DB, stream repo.DahuaStream) error {
	if stream.Internal {
		return fmt.Errorf("cannot delete internal stream")
	}

	return db.DahuaDeleteStream(ctx, stream.ID)
}

func SupportStreams(feature models.DahuaFeature) bool {
	return feature.EQ(models.DahuaFeatureCamera)
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

		args := []createInternalStreamsParams{}
		for i := 0; i < subtypes; i++ {
			arg := createInternalStreamsParams{
				Channel: int64(channelIndex + 1),
				Subtype: int64(i),
				Name:    names[i],
			}
			args = append(args, arg)
		}
		err := createInternalStreams(ctx, db, deviceID, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func RegisterStreams(bus *event.Bus, db repo.DB, store *Store) {
	bus.OnDahuaDeviceCreated(func(ctx context.Context, evt event.DahuaDeviceCreated) error {
		if SupportStreams(evt.Device.DahuaDevice.Feature) {
			// TODO: this should run on a different goroutine
			return SyncStreams(ctx, db, evt.Device.DahuaDevice.ID, store.Client(ctx, NewConn(evt.Device)).RPC)
		}
		return nil
	})
	bus.OnDahuaDeviceUpdated(func(ctx context.Context, evt event.DahuaDeviceUpdated) error {
		if SupportStreams(evt.Device.DahuaDevice.Feature) {
			// TODO: this should run on a different goroutine
			return SyncStreams(ctx, db, evt.Device.DahuaDevice.ID, store.Client(ctx, NewConn(evt.Device)).RPC)
		}
		return nil
	})
}
