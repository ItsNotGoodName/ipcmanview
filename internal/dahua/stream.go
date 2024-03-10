package dahua

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/mediamtx"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/encode"
)

func UpdateStream(ctx context.Context, stream repo.DahuaStream, arg repo.DahuaUpdateStreamParams) (repo.DahuaStream, error) {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return repo.DahuaStream{}, err
	}
	return app.DB.C().DahuaUpdateStream(ctx, arg)
}

func DeleteStream(ctx context.Context, stream repo.DahuaStream) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	if stream.Internal {
		return fmt.Errorf("cannot delete internal stream")
	}

	return app.DB.C().DahuaDeleteStream(ctx, stream.ID)
}

func SupportStream(feature models.DahuaFeature) bool {
	return feature.EQ(models.DahuaFeature_Camera)
}

// SyncStreams fetches streams from device and sync them with database.
func SyncStreams(ctx context.Context, deviceID int64, conn dahuarpc.Conn) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

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

		args := []upsertStreamsParams{}
		for i := 0; i < subtypes; i++ {
			arg := upsertStreamsParams{
				Channel: int64(channelIndex + 1),
				Subtype: int64(i),
				Name:    names[i],
			}
			args = append(args, arg)
		}
		err := upsertStreams(ctx, deviceID, args)
		if err != nil {
			return err
		}
	}

	return nil
}

type upsertStreamsParams struct {
	Channel int64
	Subtype int64
	Name    string
}

func upsertStreams(ctx context.Context, deviceID int64, arg []upsertStreamsParams) error {
	tx, err := app.DB.BeginTx(ctx, true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.C().DahuaUpdateStreamForInternal(ctx, deviceID)
	if err != nil {
		return err
	}

	ids := make([]int64, 0, len(arg))
	for _, arg := range arg {
		id, err := tx.C().DahuaCreateStreamForInternal(ctx, repo.DahuaCreateStreamForInternalParams{
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

// PushStreams pushes streams to mediamtx
func PushStreams(ctx context.Context, deviceID int64) error {
	if _, err := core.AssertAdmin(ctx); err != nil {
		return err
	}

	device, err := GetDevice(ctx, GetDeviceFilter{ID: deviceID})
	if err != nil {
		return err
	}

	streams, err := app.DB.C().DahuaListStreamsByDevice(ctx, deviceID)
	if err != nil {
		return err
	}

	for _, stream := range streams {
		name := app.MediamtxConfig.DahuaEmbedPath(stream)
		rtspURL := GetLiveRTSPURL(GetLiveRTSPURLParams{
			Username: device.Username,
			Password: device.Password,
			Host:     device.Ip,
			Port:     554,
			Channel:  int(stream.Channel),
			Subtype:  int(stream.Subtype),
		})
		rtspTransport := "tcp"
		pathConf := mediamtx.PathConf{
			Source:        &rtspURL,
			RtspTransport: &rtspTransport,
		}

		rsp, err := app.MediamtxClient.ConfigPathsGet(ctx, name)
		if err != nil {
			return err
		}
		res, err := mediamtx.ParseConfigPathsGetResponse(rsp)
		if err != nil {
			return err
		}

		switch res.StatusCode() {
		case http.StatusOK:
			rsp, err := app.MediamtxClient.ConfigPathsPatch(ctx, name, pathConf)
			if err != nil {
				return err
			}
			rsp.Body.Close()
		case http.StatusNotFound, http.StatusInternalServerError:
			rsp, err := app.MediamtxClient.ConfigPathsAdd(ctx, name, pathConf)
			if err != nil {
				return err
			}
			rsp.Body.Close()
		default:
			return errors.New(string(res.Body))
		}
	}

	return nil
}
