package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/build"
	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/coaxialcontrolio"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
)

type App struct {
	DB           *sqlx.DB
	AFS          afero.Fs
	AFSDirectory string
	DahuaStore   *dahua.Store
}

func NewConfig() huma.Config {
	return huma.DefaultConfig("IPCManView API", "1.0.0")
}

func Register(api huma.API, app App) {
	huma.Register(api, huma.Operation{
		Summary: "List devices",
		Method:  http.MethodGet,
		Path:    "/api/devices",
	}, func(ctx context.Context, input *struct{}) (*ListDevicesOutput, error) {
		body, err := ListDevices(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &ListDevicesOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*GetDeviceOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		return &GetDeviceOutput{
			Body: NewDevice(device),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Create device",
		Method:  http.MethodPost,
		Path:    "/api/devices/create",
	}, func(ctx context.Context, input *CreateDeviceInput) (*CreateDevicesOutput, error) {
		key, err := dahua.CreateDevice(ctx, app.DB, input.Body.Convert())
		if err != nil {
			return nil, err
		}

		device, err := useDevice(ctx, app.DB, key)
		if err != nil {
			return nil, err
		}

		return &CreateDevicesOutput{
			Body: NewDevice(device),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Put devices",
		Method:  http.MethodPut,
		Path:    "/api/devices",
	}, func(ctx context.Context, input *PutDevicesInput) (*ListDevicesOutput, error) {
		var args []dahua.CreateDeviceArgs
		for _, arg := range input.Body {
			args = append(args, arg.Convert())
		}

		_, err := dahua.PutDevices(ctx, app.DB, args)
		if err != nil {
			return nil, err
		}

		body, err := ListDevices(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &ListDevicesOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Update device",
		Method:  http.MethodPost,
		Path:    "/api/devices/{uuid}",
	}, func(ctx context.Context, input *struct {
		UUIDPath
		Body UpdateDevice
	},
	) (*UpdateDevicesOutput, error) {
		key, err := dahua.UpdateDevice(ctx, app.DB, dahua.UpdateDeviceArgs{
			UUID:            input.UUID,
			Name:            input.Body.Name,
			IP:              input.Body.IP,
			Username:        input.Body.Username,
			Password:        core.NullToSQLNull(input.Body.Password),
			Location:        input.Body.Location,
			Features:        types.NewSlice(input.Body.Features),
			Email:           input.Body.Email,
			Latitude:        core.NullToSQLNull(input.Body.Latitude),
			Longitude:       core.NullToSQLNull(input.Body.Longitude),
			SunriseOffset:   input.Body.SunriseOffset,
			SunsetOffset:    input.Body.SunsetOffset,
			SyncVideoInMode: core.NullToSQLNull(input.Body.SyncVideoInMode),
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, huma.Error404NotFound("device not found")
			}
			return nil, err
		}

		device, err := useDevice(ctx, app.DB, key)
		if err != nil {
			return nil, err
		}

		return &UpdateDevicesOutput{
			Body: NewDevice(device),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Delete device",
		Method:  http.MethodDelete,
		Path:    "/api/devices/{uuid}",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*struct{}, error) {
		return &struct{}{}, dahua.DeleteDevice(ctx, app.DB, input.UUID)
	})
	huma.Register(api, huma.Operation{
		Summary: "Get home page",
		Method:  http.MethodGet,
		Path:    "/api/pages/home",
	}, func(ctx context.Context, input *struct{}) (*GetHomePageOutput, error) {
		fileUsage, err := core.DirectorySize(app.AFSDirectory)
		if err != nil {
			return nil, err
		}

		body := GetHomePage{
			FileUsage: fileUsage,
			Build:     build.Current,
		}
		err = app.DB.GetContext(ctx, &body, `
			SELECT 
				(SELECT count(*) FROM dahua_devices) AS device_count,
				(SELECT count(*) FROM dahua_events) AS event_count,
				(SELECT count(*) FROM dahua_email_messages) AS email_count,
				(SELECT count(*) FROM dahua_files) AS file_count,
				(SELECT page_count * page_size as size FROM pragma_page_count(), pragma_page_size()) AS db_usage
		`)
		if err != nil {
			return nil, err
		}

		return &GetHomePageOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device coaxial caps",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/coaxial/caps",
	}, func(ctx context.Context, input *struct {
		UUIDPath
		Channel int `query:"channel"`
	},
	) (*GetDeviceCoaxialCapsOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.GetCoaxialCaps(ctx, client.RPC, input.Channel)
		if err != nil {
			return nil, err
		}

		return &GetDeviceCoaxialCapsOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device coaxial status",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/coaxial/status",
	}, func(ctx context.Context, input *struct {
		UUIDPath
		Channel int `query:"channel"`
	},
	) (*GetDeviceCoaxialStatusOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.GetCoaxialStatus(ctx, client.RPC, input.Channel)
		if err != nil {
			return nil, err
		}

		return &GetDeviceCoaxialStatusOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device detail",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/detail",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*GetDeviceDetailOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.GetDetail(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &GetDeviceDetailOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device licenses",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/licenses",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*GetDeviceLicensesOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ListLicenses(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &GetDeviceLicensesOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List device ptz presets",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/ptz/presets",
	}, func(ctx context.Context, input *struct {
		UUIDPath
		Channel int `query:"channel"`
	},
	) (*ListDevicePTZPresetsOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ListPTZPresets(ctx, client.PTZ, input.Channel)
		if err != nil {
			return nil, err
		}

		return &ListDevicePTZPresetsOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device software versions",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/software",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*GetDeviceSoftwareOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.GetSoftwareVersion(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &GetDeviceSoftwareOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List device storage",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/storage",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*ListDeviceStorageOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ListStorage(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &ListDeviceStorageOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List device users",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/users",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*ListDeviceUsersOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		position, err := dahua.GetDevicePosition(ctx, app.DB, device.ID)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ListUsers(ctx, client.RPC, position.Location.Location)
		if err != nil {
			return nil, err
		}

		return &ListDeviceUsersOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List device active users",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/users/active",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*ListDeviceActiveUsersOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		position, err := dahua.GetDevicePosition(ctx, app.DB, device.ID)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ListActiveUsers(ctx, client.RPC, position.Location.Location)
		if err != nil {
			return nil, err
		}

		return &ListDeviceActiveUsersOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List device groups",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/groups",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*ListDeviceGroupsOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ListGroups(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &ListDeviceGroupsOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device uptime",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/uptime",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*GetDeviceUptimeOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.GetUptime(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &GetDeviceUptimeOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device status",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/status",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*GetDeviceStatusOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body := dahua.GetStatus(ctx, client.RPC)

		return &GetDeviceStatusOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device snapshot",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/snapshot",
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Current snapshot of camera",
				Content:     map[string]*huma.MediaType{"image/jpeg": {}},
			},
		},
	}, func(ctx context.Context, input *struct {
		UUIDPath
		Channel int `query:"channel"`
		Type    int `query:"type"`
	},
	) (*huma.StreamResponse, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		snapshot, err := dahuacgi.SnapshotGet(ctx, client.CGI, input.Channel, input.Type)
		if err != nil {
			return nil, err
		}

		return &huma.StreamResponse{
			Body: func(ctx huma.Context) {
				defer snapshot.Close()

				ctx.SetHeader("Content-Type", snapshot.ContentType)
				ctx.SetHeader("Content-Length", snapshot.ContentLength)

				io.Copy(ctx.BodyWriter(), snapshot)
			},
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List event codes",
		Method:  http.MethodGet,
		Path:    "/api/event-codes",
	}, func(ctx context.Context, input *struct{},
	) (*EventCodesOutput, error) {
		body, err := dahua.ListEventCodes(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &EventCodesOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List event actions",
		Method:  http.MethodGet,
		Path:    "/api/event-actions",
	}, func(ctx context.Context, input *struct{},
	) (*EventActionsOutput, error) {
		body, err := dahua.ListEventActions(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &EventActionsOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Delete events",
		Method:  http.MethodDelete,
		Path:    "/api/events",
	}, func(ctx context.Context, input *struct{},
	) (*struct{}, error) {
		err := dahua.DeleteEvents(ctx, app.DB)
		if err != nil {
			return nil, err
		}
		return &struct{}{}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List events",
		Method:  http.MethodGet,
		Path:    "/api/events",
	}, func(ctx context.Context, input *struct {
		PageQuery
		OrderQuery
		Devices []string `query:"device"`
		Codes   []string `query:"codes"`
		Actions []string `query:"actions"`
	},
	) (*ListEventsOutput, error) {
		deviceIDs, err := dahua.GetDeviceIDs(ctx, app.DB, input.Devices)
		if err != nil {
			return nil, err
		}

		res, err := dahua.ListEvents(ctx, app.DB, dahua.ListEventsParams{
			Page:      input.GetPage(),
			Ascending: input.Order.Ascending(),
			Filter: dahua.ListEventsFilter{
				DeviceIDs: deviceIDs,
				Codes:     input.Codes,
				Actions:   input.Actions,
			},
		})
		if err != nil {
			return nil, err
		}

		data := []DeviceEvent{}
		for _, v := range res.Items {
			data = append(data, DeviceEvent{
				ID:         v.ID,
				DeviceUUID: v.Device_UUID,
				Code:       v.Code,
				Action:     v.Action,
				Index:      v.Index,
				Data:       v.Data.RawMessage,
				CreatedAt:  v.Created_At.Time,
			})
		}

		return &ListEventsOutput{
			Body: ListEvents{
				Pagination: NewPagePagination(res.PageResult),
				Data:       data,
			},
		}, nil
	})
	sse.Register(api, huma.Operation{
		Summary: "Listen for events",
		Method:  http.MethodGet,
		Path:    "/api/events/sse",
	}, map[string]any{
		"message": DeviceEvent{},
	}, func(ctx context.Context, input *struct {
		Devices []string `query:"device"`
		Codes   []string `query:"codes"`
		Actions []string `query:"actions"`
	}, send sse.Sender,
	) {
		eventC, unsub := bus.SubscribeChannel[bus.EventCreated]("app(/api/events)")
		defer unsub()

		for event := range eventC {
			if !event.AllowLive {
				continue
			}
			if len(input.Devices) != 0 && !slices.Contains(input.Devices, event.DeviceKey.UUID) {
				continue
			}
			if len(input.Codes) != 0 && !slices.Contains(input.Codes, event.Event.Code) {
				continue
			}
			if len(input.Actions) != 0 && !slices.Contains(input.Actions, event.Event.Action) {
				continue
			}
			err := send.Data(DeviceEvent{
				ID:         event.EventID,
				DeviceUUID: event.DeviceKey.UUID,
				Code:       event.Event.Code,
				Action:     event.Event.Action,
				Index:      int64(event.Event.Index),
				Data:       event.Event.Data,
				CreatedAt:  event.CreatedAt,
			})
			if err != nil {
				return
			}
		}
	})
	huma.Register(api, huma.Operation{
		Summary: "Download device file",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/file",
		Responses: map[string]*huma.Response{
			"200": {
				Description: "File from camera",
				Content: map[string]*huma.MediaType{
					"image/jpeg":               {},
					"application/octet-stream": {},
				},
			},
		},
	}, func(ctx context.Context, input *struct {
		UUIDPath
		// TODO: this should be path param wildcard but OpenAPI is stupid
		Name string `query:"name" required:"true"`
	},
	) (*huma.StreamResponse, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		var file dahua.File
		err = app.DB.GetContext(ctx, &file, `
			SELECT * FROM dahua_files WHERE file_path = ?
		`, input.Name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, huma.Error404NotFound("file not found")
			}
			return nil, err
		}

		var (
			rd            io.ReadCloser
			contentLength int64
		)
		switch file.Storage {
		case dahua.StorageLocal:
			rd, contentLength, err = dahua.OpenFileLocal(ctx, client, file.File_Path)
			if err != nil {
				return nil, err
			}
		case dahua.StorageFTP:
			rd, contentLength, err = dahua.OpenFileFTP(ctx, app.DB, file.File_Path)
			if err != nil {
				return nil, err
			}
		case dahua.StorageSFTP:
			rd, contentLength, err = dahua.OpenFileSFTP(ctx, app.DB, file.File_Path)
			if err != nil {
				return nil, err
			}
		default:
			return nil, huma.Error500InternalServerError(fmt.Sprintf("unknown storage: %s", file.Storage))
		}

		return &huma.StreamResponse{
			Body: func(ctx huma.Context) {
				defer rd.Close()
				ctx.SetHeader("Content-Length", strconv.FormatInt(contentLength, 10))
				ctx.SetHeader("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filepath.Base(input.Name)))
				io.Copy(ctx.BodyWriter(), rd)
			},
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Reboot device",
		Method:  http.MethodPost,
		Path:    "/api/devices/{uuid}/reboot",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*struct{}, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		return &struct{}{}, dahua.RebootDevice(ctx, client.RPC)
	})
	huma.Register(api, huma.Operation{
		Summary: "Get device VideoInMode",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/video-in-mode",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*DeviceVideoInModeOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.GetVideoInMode(ctx, client.RPC)
		if err != nil {
			return nil, err
		}

		return &DeviceVideoInModeOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Reset device file scan cursor",
		Method:  http.MethodDelete,
		Path:    "/api/devices/{uuid}/scan/cursor",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*DeviceScanCursorOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		if err := dahua.ResetScanCursor(ctx, app.DB, device.ID); err != nil {
			return nil, err
		}

		body, err := GetFileCursor(ctx, app.DB, device.ID)
		if err != nil {
			return nil, err
		}

		return &DeviceScanCursorOutput{
			Body: NewFileScanCursor(body),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get file scan cursor",
		Method:  http.MethodGet,
		Path:    "/api/devices/{uuid}/scan/cursor",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*DeviceScanCursorOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		body, err := GetFileCursor(ctx, app.DB, device.ID)
		if err != nil {
			return nil, err
		}

		return &DeviceScanCursorOutput{
			Body: NewFileScanCursor(body),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Manual scan device files",
		Method:  http.MethodPost,
		Path:    "/api/devices/{uuid}/scan/manual",
	}, func(ctx context.Context, input *struct {
		UUIDPath
		Body ManualFileScan
	},
	) (*DeviceScanOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		start := core.Optional(input.Body.StartTime, dahua.ScanEpoch)
		end := core.Optional(input.Body.EndTime, time.Now())

		body, err := dahua.ScanManual(ctx, app.DB, client.RPC, device.ID, start, end)
		if err != nil {
			return nil, err
		}

		return &DeviceScanOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Full scan device files",
		Method:  http.MethodPost,
		Path:    "/api/devices/{uuid}/scan/full",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*DeviceScanOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ScanFull(ctx, app.DB, client.RPC, device.ID)
		if err != nil {
			return nil, err
		}

		return &DeviceScanOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Quick scan device files",
		Method:  http.MethodPost,
		Path:    "/api/devices/{uuid}/scan/quick",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*DeviceScanOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		body, err := dahua.ScanQuick(ctx, app.DB, client.RPC, device.ID)
		if err != nil {
			return nil, err
		}

		return &DeviceScanOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Sync device VideoInMode",
		Method:  http.MethodPost,
		Path:    "/api/devices/{uuid}/video-in-mode/sync",
	}, func(ctx context.Context, input *struct {
		UUIDPath
		Body DeviceVideoInModeSync
	},
	) (*DeviceVideoInModeOutput, error) {
		device, err := useDevice(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		client, err := useClient(ctx, app.DahuaStore, device)
		if err != nil {
			return nil, err
		}

		position, err := dahua.GetDevicePosition(ctx, app.DB, device.ID)
		if err != nil {
			return nil, err
		}

		body, err := dahua.SyncVideoInMode(ctx, client.RPC, dahua.SyncVideoInModeArgs{
			Location:      core.Optional(input.Body.Location, position.Location).Location,
			Latitude:      core.Optional(input.Body.Latitude, position.Latitude),
			Longitude:     core.Optional(input.Body.Longitude, position.Longitude),
			SunriseOffset: core.Optional(input.Body.SunriseOffset, position.Sunrise_Offset).Duration,
			SunsetOffset:  core.Optional(input.Body.SunsetOffset, position.Sunset_Offset).Duration,
		})
		if err != nil {
			return nil, err
		}

		return &DeviceVideoInModeOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Set device white light state",
		Method:  http.MethodPost,
		Path:    "/api/devices/{uuid}/coaxial/white-light",
	}, func(ctx context.Context, input *struct {
		UUIDPath
		Channel int    `query:"channel"`
		Action  string `query:"action" enum:"on,off,toggle"`
	},
	) (*struct{}, error) {
		return &struct{}{}, SetDeviceCoaxialState(ctx, app.DahuaStore, app.DB, input.UUID, input.Channel, coaxialcontrolio.TypeWhiteLight, input.Action)
	})
	huma.Register(api, huma.Operation{
		Summary: "Set device speaker state",
		Method:  http.MethodPost,
		Path:    "/api/devices/{uuid}/coaxial/speaker",
	}, func(ctx context.Context, input *struct {
		UUIDPath
		Channel int    `query:"channel"`
		Action  string `query:"action" enum:"on,off,toggle"`
	},
	) (*struct{}, error) {
		return &struct{}{}, SetDeviceCoaxialState(ctx, app.DahuaStore, app.DB, input.UUID, input.Channel, coaxialcontrolio.TypeSpeaker, input.Action)
	})
	huma.Register(api, huma.Operation{
		Summary: "Get settings",
		Method:  http.MethodGet,
		Path:    "/api/settings",
	}, func(ctx context.Context, i *struct{}) (*SettingOutput, error) {
		settings, err := system.GetSettings(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &SettingOutput{
			Body: NewSettings(settings),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Update settings",
		Method:  http.MethodPut,
		Path:    "/api/settings",
	}, func(ctx context.Context, input *struct {
		Body UpdateSettings
	},
	) (*SettingOutput, error) {
		err := system.UpdateSettings(ctx, app.DB, system.UpdateSettingsArgs{
			Location:        input.Body.Location,
			Latitude:        input.Body.Latitude,
			Longitude:       input.Body.Longitude,
			SunriseOffset:   input.Body.SunriseOffset,
			SunsetOffset:    input.Body.SunsetOffset,
			SyncVideoInMode: input.Body.SyncVideoInMode,
		})
		if err != nil {
			return nil, err
		}

		settings, err := system.GetSettings(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &SettingOutput{
			Body: NewSettings(settings),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Patch settings",
		Method:  http.MethodPatch,
		Path:    "/api/settings",
	}, func(ctx context.Context, input *struct {
		Body PatchSettings
	},
	) (*SettingOutput, error) {
		err := system.PatchSettings(ctx, app.DB, system.PatchSettingsArgs{
			Location:        input.Body.Location,
			Latitude:        input.Body.Latitude,
			Longitude:       input.Body.Longitude,
			SunriseOffset:   input.Body.SunriseOffset,
			SunsetOffset:    input.Body.SunsetOffset,
			SyncVideoInMode: input.Body.SyncVideoInMode,
		})
		if err != nil {
			return nil, err
		}

		settings, err := system.GetSettings(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &SettingOutput{
			Body: NewSettings(settings),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Default settings",
		Method:  http.MethodDelete,
		Path:    "/api/settings",
	}, func(ctx context.Context, i *struct{}) (*SettingOutput, error) {
		err := system.DefaultSettings(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		settings, err := system.GetSettings(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &SettingOutput{
			Body: NewSettings(settings),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List email endpoints",
		Method:  http.MethodGet,
		Path:    "/api/email-endpoints",
	}, func(ctx context.Context, i *struct{}) (*ListEmailEndpointsOutput, error) {
		body, err := ListEmailEndpoints(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &ListEmailEndpointsOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Get email endpoint",
		Method:  http.MethodGet,
		Path:    "/api/email-endpoints/{uuid}",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*GetEmailEndpointOutput, error) {
		body, err := GetEmailEndpoints(ctx, app.DB, types.Key{UUID: input.UUID})
		if err != nil {
			return nil, err
		}

		return &GetEmailEndpointOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Create email endpoint",
		Method:  http.MethodPost,
		Path:    "/api/email-endpoints/create",
	}, func(ctx context.Context, input *struct {
		Body CreateEmailEndpoint
	},
	) (*CreateEmailEndpointOutput, error) {
		key, err := dahua.CreateEmailEndpoint(ctx, app.DB, input.Body.Convert())
		if err != nil {
			return nil, err
		}

		body, err := GetEmailEndpoints(ctx, app.DB, key)
		if err != nil {
			return nil, err
		}

		return &CreateEmailEndpointOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Put email endpoints",
		Method:  http.MethodPut,
		Path:    "/api/email-endpoints",
	}, func(ctx context.Context, input *struct {
		Body []CreateEmailEndpoint
	},
	) (*ListEmailEndpointsOutput, error) {
		var args []dahua.CreateEmailEndpointArgs
		for _, v := range input.Body {
			args = append(args, v.Convert())
		}

		_, err := dahua.PutEmailEndpoints(ctx, app.DB, args)
		if err != nil {
			return nil, err
		}

		body, err := ListEmailEndpoints(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &ListEmailEndpointsOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Update email endpoint",
		Method:  http.MethodPost,
		Path:    "/api/email-endpoints/{uuid}",
	}, func(ctx context.Context, input *struct {
		UUIDPath
		Body UpdateEmailEndpoint
	},
	) (*GetEmailEndpointOutput, error) {
		key, err := dahua.UpdateEmailEndpoint(ctx, app.DB, dahua.UpdateEmailEndpointArgs{
			UUID:          input.UUID,
			Global:        input.Body.Global,
			Expression:    input.Body.Expression,
			TitleTemplate: input.Body.TitleTemplate,
			BodyTemplate:  input.Body.BodyTemplate,
			Attachments:   input.Body.Attachments,
			URLs:          types.NewSlice(input.Body.URLs),
			DeviceUUIDs:   input.Body.DeviceUUIDs,
			Disabled:      input.Body.Disabled,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, huma.Error404NotFound("email endpoint not found")
			}
			return nil, err
		}

		body, err := GetEmailEndpoints(ctx, app.DB, key)
		if err != nil {
			return nil, err
		}

		return &GetEmailEndpointOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Delete email endpoint",
		Method:  http.MethodDelete,
		Path:    "/api/email-endpoints/{uuid}",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*struct{}, error) {
		return &struct{}{}, dahua.DeleteEmailEndpoint(ctx, app.DB, input.UUID)
	})
	huma.Register(api, huma.Operation{
		Summary: "List storage destinations",
		Method:  http.MethodGet,
		Path:    "/api/storage-destinations",
	}, func(ctx context.Context, input *struct{}) (*ListStorageDestinationOutput, error) {
		body, err := ListStorageDestinations(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &ListStorageDestinationOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Put storage destinations",
		Method:  http.MethodPut,
		Path:    "/api/storage-destinations",
	}, func(ctx context.Context, input *PutCreateStorageDestinationInput) (*ListStorageDestinationOutput, error) {
		var args []dahua.CreateStorageDestinationArgs
		for _, arg := range input.Body {
			args = append(args, arg.Convert())
		}

		_, err := dahua.PutStorageDestinations(ctx, app.DB, args)
		if err != nil {
			return nil, err
		}

		body, err := ListStorageDestinations(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &ListStorageDestinationOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List locations",
		Method:  http.MethodGet,
		Path:    "/api/locations",
	}, func(ctx context.Context, input *struct{},
	) (*LocationsOutput, error) {
		return &LocationsOutput{
			Body: []string{
				"Local",
				"UTC",
			},
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List event rules",
		Method:  http.MethodGet,
		Path:    "/api/event-rules",
	}, func(ctx context.Context, input *struct{},
	) (*ListEventRulesOutput, error) {
		body, err := ListEventRules(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &ListEventRulesOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Create event rule",
		Method:  http.MethodPost,
		Path:    "/api/event-rules/create",
	}, func(ctx context.Context, input *struct {
		Body CreateEventRule
	},
	) (*CreateEventRulesOutput, error) {
		body, err := dahua.CreateEventRule(ctx, app.DB, dahua.CreateEventRuleArgs{
			UUID:      core.Optional(input.Body.UUID, uuid.NewString()),
			Code:      input.Body.Code,
			AllowDB:   input.Body.AllowDB,
			AllowLive: input.Body.AllowLive,
			AllowMQTT: input.Body.AllowMQTT,
		})
		if err != nil {
			return nil, err
		}

		return &CreateEventRulesOutput{
			Body: NewEventRule(body),
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Update event rules",
		Method:  http.MethodPost,
		Path:    "/api/event-rules/{uuid}",
	}, func(ctx context.Context, input *struct {
		UUIDPath
		Body UpdateEventRule
	},
	) (*ListEventRulesOutput, error) {
		err := dahua.UpdateEventRule(ctx, app.DB, dahua.UpdateEventRuleArgs{
			UUID:      input.UUID,
			Code:      input.Body.Code,
			AllowDB:   input.Body.AllowDB,
			AllowLive: input.Body.AllowLive,
			AllowMQTT: input.Body.AllowMQTT,
		})
		if err != nil {
			return nil, err
		}

		body, err := ListEventRules(ctx, app.DB)
		if err != nil {
			return nil, err
		}

		return &ListEventRulesOutput{
			Body: body,
		}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Delete event rule by uuid",
		Method:  http.MethodDelete,
		Path:    "/api/event-rules/{uuid}",
	}, func(ctx context.Context, input *struct {
		UUIDPath
	},
	) (*struct{}, error) {
		err := dahua.DeleteEventRule(ctx, app.DB, input.UUID)
		if err != nil {
			return nil, err
		}
		return &struct{}{}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "Delete event rules",
		Method:  http.MethodDelete,
		Path:    "/api/event-rules",
	}, func(ctx context.Context, input *struct {
		Body []string
	},
	) (*struct{}, error) {
		for _, uuid := range input.Body {
			err := dahua.DeleteEventRule(ctx, app.DB, uuid)
			if err != nil {
				return nil, err
			}
		}
		return &struct{}{}, nil
	})
	huma.Register(api, huma.Operation{
		Summary: "List files",
		Method:  http.MethodGet,
		Path:    "/api/files",
	}, func(ctx context.Context, input *struct {
		PageQuery
		Devices []string `query:"device"`
	},
	) (*ListFilesOutput, error) {
		return &ListFilesOutput{
			Body: ListFiles{
				Pagination: PagePagination{
					Page: 1,
				},
				Data: []File{},
			},
		}, nil
	})
}

func NewEmailEndpoint(v dahua.EmailEndpoint, deviceUUIDs []string) EmailEndpoint {
	return EmailEndpoint{
		UUID:          v.UUID,
		Global:        v.Global,
		DeviceUUIDs:   deviceUUIDs,
		Expression:    v.Expression,
		URLs:          v.URLs.V,
		TitleTemplate: v.Title_Template,
		BodyTemplate:  v.Body_Template,
		Attachments:   v.Attachments,
		CreatedAt:     v.Created_At.Time,
		UpdatedAt:     v.Updated_At.Time,
		Disabled:      v.Disabled_At.Valid,
	}
}

func NewSettings(v system.Settings) Settings {
	return Settings{
		Location:        v.Location,
		Latitude:        v.Latitude,
		Longitude:       v.Longitude,
		SunriseOffset:   v.Sunrise_Offset,
		SunsetOffset:    v.Sunset_Offset,
		SyncVideoInMode: v.Sync_Video_In_Mode,
	}
}

type SettingOutput struct {
	Body Settings
}

type DeviceVideoInModeOutput struct {
	Body dahua.DeviceVideoInMode
}

type ListDevicesOutput struct {
	Body []Device
}

func NewDevice(v dahua.Device) Device {
	ip := net.ParseIP(v.IP)
	return Device{
		UUID:            v.UUID,
		Name:            v.Name,
		IP:              ip,
		Username:        v.Username,
		Location:        core.SQLNullToNull(v.Location),
		Features:        v.Features.V,
		Email:           v.Email.String,
		CreatedAt:       v.Created_At.Time,
		UpdatedAt:       v.Updated_At.Time,
		Latitude:        core.SQLNullToNull(v.Latitude),
		Longitude:       core.SQLNullToNull(v.Longitude),
		SunriseOffset:   core.SQLNullToNull(v.Sunrise_Offset),
		SunsetOffset:    core.SQLNullToNull(v.Sunset_Offset),
		SyncVideoInMode: core.SQLNullToNull(v.Sync_Video_In_Mode),
	}
}

func (i *CreateDevice) Resolve(ctx huma.Context) []error {
	if i.Name == "" {
		i.Name = i.IP
	}
	return nil
}

func (i *CreateDevice) Convert() dahua.CreateDeviceArgs {
	return dahua.CreateDeviceArgs{
		UUID:            core.Optional(i.UUID, uuid.NewString()),
		Name:            i.Name,
		IP:              i.IP,
		Username:        i.Username,
		Password:        i.Password,
		Location:        i.Location,
		Features:        types.NewSlice(i.Features),
		Email:           core.NullToSQLNull(i.Email),
		Latitude:        core.NullToSQLNull(i.Latitude),
		Longitude:       core.NullToSQLNull(i.Longitude),
		SunriseOffset:   i.SunriseOffset,
		SunsetOffset:    i.SunsetOffset,
		SyncVideoInMode: core.NullToSQLNull(i.SyncVideoInMode),
	}
}

type CreateDeviceInput struct {
	Body CreateDevice
}

type PutDevicesInput struct {
	Body []CreateDevice
}

type PutDevicesOutput struct {
	Body []Device
}

type CreateDevicesOutput struct {
	Body Device
}

type GetDeviceOutput struct {
	Body Device
}

type UpdateDevicesOutput struct {
	Body Device
}

type PatchDevicesOutput struct {
	Body Device
}

type GetDeviceCoaxialCapsOutput struct {
	Body dahua.DeviceCoaxialCaps
}

type GetDeviceCoaxialStatusOutput struct {
	Body dahua.DeviceCoaxialStatus
}

type GetDeviceDetailOutput struct {
	Body dahua.DeviceDetail
}

type GetDeviceLicensesOutput struct {
	Body []dahua.DeviceLicense
}

type ListDevicePTZPresetsOutput struct {
	Body []dahua.DevicePTZPreset
}

type GetDeviceSoftwareOutput struct {
	Body dahua.DeviceSoftwareVersion
}

type ListDeviceStorageOutput struct {
	Body []dahua.DeviceStorage
}

type ListDeviceUsersOutput struct {
	Body []dahua.DeviceUser
}

type ListDeviceActiveUsersOutput struct {
	Body []dahua.DeviceActiveUser
}

type ListDeviceGroupsOutput struct {
	Body []dahua.DeviceGroup
}

type GetDeviceUptimeOutput struct {
	Body dahua.DeviceUptime
}

type GetDeviceStatusOutput struct {
	Body dahua.DeviceStatus
}

type GetDeviceSnapshotOutput struct {
	Body []byte
}

func (i CreateEmailEndpoint) Convert() dahua.CreateEmailEndpointArgs {
	return dahua.CreateEmailEndpointArgs{
		UUID:          core.Optional(i.UUID, uuid.NewString()),
		Global:        i.Global,
		Expression:    i.Expression,
		TitleTemplate: core.Optional(i.TitleTemplate, "{{.Message.Subject}}"),
		BodyTemplate:  core.Optional(i.BodyTemplate, "{{.Message.Text}}"),
		Attachments:   i.Attachments,
		URLs:          types.NewSlice(i.URLs),
		DeviceUUIDs:   i.DeviceUUIDs,
		Disabled:      i.Disabled,
	}
}

type ListEmailEndpointsOutput struct {
	Body []EmailEndpoint
}

type GetEmailEndpointOutput struct {
	Body EmailEndpoint
}

type CreateEmailEndpointOutput struct {
	Body EmailEndpoint
}

type GetHomePageOutput struct {
	Body GetHomePage
}

func ListEmailEndpoints(ctx context.Context, db *sqlx.DB) ([]EmailEndpoint, error) {
	rows, err := db.QueryxContext(ctx, `
		SELECT * FROM dahua_email_endpoints
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	body := []EmailEndpoint{}
	for rows.Next() {
		var v dahua.EmailEndpoint
		if err := rows.StructScan(&v); err != nil {
			return nil, err
		}

		deviceUUIDs, err := dahua.GetEmailEndpointDeviceUUIDs(ctx, db, v.Key)
		if err != nil {
			return nil, err
		}

		body = append(body, NewEmailEndpoint(v, deviceUUIDs))
	}

	return body, nil
}

func GetEmailEndpoints(ctx context.Context, db *sqlx.DB, key types.Key) (EmailEndpoint, error) {
	var emailEndpoint dahua.EmailEndpoint
	err := db.GetContext(ctx, &emailEndpoint, `
		SELECT * FROM dahua_email_endpoints WHERE id = ? OR uuid = ? LIMIT 1;
	`, key.ID, key.UUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return EmailEndpoint{}, huma.Error404NotFound("email endpoint not found")
		}
		return EmailEndpoint{}, err
	}

	deviceUUIDs, err := dahua.GetEmailEndpointDeviceUUIDs(ctx, db, emailEndpoint.Key)
	if err != nil {
		return EmailEndpoint{}, err
	}

	return NewEmailEndpoint(emailEndpoint, deviceUUIDs), nil
}

func ListDevices(ctx context.Context, db *sqlx.DB) ([]Device, error) {
	rows, err := db.QueryxContext(ctx, `
		SELECT * FROM dahua_devices
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	body := []Device{}
	for rows.Next() {
		var v dahua.Device
		if err := rows.StructScan(&v); err != nil {
			return nil, err
		}
		body = append(body, NewDevice(v))
	}

	return body, nil
}

func SetDeviceCoaxialState(ctx context.Context, dahuaStore *dahua.Store, db *sqlx.DB, uuid string, channel int, typE coaxialcontrolio.Type, action string) error {
	device, err := useDevice(ctx, db, types.Key{UUID: uuid})
	if err != nil {
		return err
	}

	client, err := useClient(ctx, dahuaStore, device)
	if err != nil {
		return err
	}

	control := coaxialcontrolio.ControlRequest{
		Type:        typE,
		IO:          coaxialcontrolio.Off,
		TriggerMode: coaxialcontrolio.TriggerModeManual,
	}
	switch action {
	case "on":
		control.IO = coaxialcontrolio.On
	case "off":
		control.IO = coaxialcontrolio.Off
	default:
		status, err := dahua.GetCoaxialStatus(ctx, client.RPC, channel)
		if err != nil {
			return err
		}
		if status.WhiteLight {
			control.IO = coaxialcontrolio.Off
		} else {
			control.IO = coaxialcontrolio.On
		}
	}

	return coaxialcontrolio.Control(ctx, client.RPC, channel, control)
}

func useDevice(ctx context.Context, db *sqlx.DB, key types.Key) (dahua.Device, error) {
	var device dahua.Device
	err := db.GetContext(ctx, &device, `
		SELECT * FROM dahua_devices WHERE id = ? OR uuid = ? LIMIT 1
	`, key.ID, key.UUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dahua.Device{}, huma.Error404NotFound("device not found")
		}
		return dahua.Device{}, err
	}
	return device, nil
}

func useClient(ctx context.Context, dahuaStore *dahua.Store, device dahua.Device) (dahua.StoreClient, error) {
	client, err := dahuaStore.GetClient(ctx, device.Key)
	if err != nil {
		return dahua.StoreClient{}, huma.Error404NotFound("device not found")
	}
	return client, nil
}

func NewStorageDestination(v dahua.StorageDestination) StorageDestination {
	return StorageDestination{
		UUID:            v.UUID,
		Name:            v.Name,
		Storage:         v.Storage,
		ServerAddress:   v.Server_Address,
		Port:            v.Port,
		Username:        v.Username,
		Password:        v.Password,
		RemoteDirectory: v.Remote_Directory,
		CreatedAt:       v.Created_At.Time,
		UpdatedAt:       v.Updated_At.Time,
	}
}

func (i *CreateStorageDestination) Convert() dahua.CreateStorageDestinationArgs {
	return dahua.CreateStorageDestinationArgs{
		UUID:            core.Optional(i.UUID, uuid.NewString()),
		Name:            i.Name,
		Storage:         i.Storage,
		ServerAddress:   i.ServerAddress,
		Port:            i.Port,
		Username:        i.Username,
		Password:        i.Password,
		RemoteDirectory: i.RemoteDirectory,
	}
}

type PutCreateStorageDestinationInput struct {
	Body []CreateStorageDestination
}

func ListStorageDestinations(ctx context.Context, db *sqlx.DB) ([]StorageDestination, error) {
	rows, err := db.QueryxContext(ctx, `
		SELECT * FROM dahua_storage_destinations
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	body := []StorageDestination{}
	for rows.Next() {
		var v dahua.StorageDestination
		if err := rows.StructScan(&v); err != nil {
			return nil, err
		}
		body = append(body, NewStorageDestination(v))
	}

	return body, nil
}

type ListStorageDestinationOutput struct {
	Body []StorageDestination
}

type DeviceScanOutput struct {
	Body dahua.ScanResult
}

type DeviceScanCursorOutput struct {
	Body FileScanCursor
}

func NewFileScanCursor(v dahua.ScanCursor) FileScanCursor {
	return FileScanCursor{
		QuickCursor:  v.Quick_Cursor.Time,
		FullCursor:   v.Full_Cursor.Time,
		FullEpoch:    v.Full_Epoch.Time,
		FullComplete: v.Full_Complete,
		UpdatedAt:    v.Updated_At.Time,
	}
}

func GetFileCursor(ctx context.Context, db *sqlx.DB, deviceID int64) (dahua.ScanCursor, error) {
	var v dahua.ScanCursor
	err := db.GetContext(ctx, &v, `
		SELECT * FROM dahua_scan_cursors WHERE device_id = ?
	`, deviceID)
	return v, err
}

type EventCodesOutput struct {
	Body []string
}

type EventActionsOutput struct {
	Body []string
}

type LocationsOutput struct {
	Body []string
}

type ListEventsOutput struct {
	Body ListEvents
}

type ListEvents struct {
	Pagination PagePagination `json:"pagination"`
	Data       []DeviceEvent  `json:"data"`
}

type CreateEventRulesOutput struct {
	Body EventRule
}

type ListEventRulesOutput struct {
	Body []EventRule
}

func ListEventRules(ctx context.Context, db *sqlx.DB) ([]EventRule, error) {
	rows, err := db.QueryxContext(ctx, `
		SELECT * FROM dahua_event_rules
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []EventRule{}
	for rows.Next() {
		var v dahua.EventRule
		if err := rows.StructScan(&v); err != nil {
			return nil, err
		}
		res = append(res, NewEventRule(v))
	}

	return res, nil
}

func NewEventRule(v dahua.EventRule) EventRule {
	return EventRule{
		UUID:      v.UUID,
		Code:      v.Code,
		AllowDB:   v.Allow_DB,
		AllowLive: v.Allow_Live,
		AllowMQTT: v.Allow_MQTT,
		CanDelete: v.Can_Delete,
	}
}

type EventRule struct {
	UUID      string `json:"uuid"`
	Code      string `json:"code"`
	AllowDB   bool   `json:"allow_db"`
	AllowLive bool   `json:"allow_live"`
	AllowMQTT bool   `json:"allow_mqtt"`
	CanDelete bool   `json:"can_delete"`
}

type ListFilesOutput struct {
	Body ListFiles
}

type ListFiles struct {
	Pagination PagePagination `json:"pagination"`
	Data       []File         `json:"data"`
}
