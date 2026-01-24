package bus

import (
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
)

type DeviceCreated struct {
	DeviceKey types.Key
}

type DeviceUpdated struct {
	DeviceKey types.Key
}

type DeviceDeleted struct {
	DeviceKeys []types.Key
}

type EventCreated struct {
	EventID   string
	DeviceKey types.Key
	AllowDB   bool
	AllowMQTT bool
	AllowLive bool
	Event     dahuacgi.Event
	CreatedAt time.Time
}

type CoaxialStatusUpdated struct {
	DeviceKey  types.Key
	Channel    int
	WhiteLight bool
	Speaker    bool
}

type EmailCreated struct {
	DeviceKey  types.Key
	MessageKey types.Key
}

type FileScanProgressed struct {
	DeviceKey types.Key
	Progress  float64
}
