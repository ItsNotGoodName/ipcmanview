package action

import "github.com/ItsNotGoodName/ipcmanview/internal/event"

var (
	DahuaDeviceCreated = event.NewEventBuilder[int64]("dahua-device:created")
	DahuaDeviceUpdated = event.NewEventBuilder[int64]("dahua-device:updated")
	DahuaDeviceDeleted = event.NewEventBuilder[int64]("dahua-device:deleted")
	DahuaEmailCreated  = event.NewEventBuilder[int64]("dahua-email:created")
)
