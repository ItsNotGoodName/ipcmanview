package action

import "github.com/ItsNotGoodName/ipcmanview/internal/system"

var (
	DahuaDeviceCreated = system.NewEventBuilder[int64]("dahua-device:created")
	DahuaDeviceUpdated = system.NewEventBuilder[int64]("dahua-device:updated")
	DahuaDeviceDeleted = system.NewEventBuilder[int64]("dahua-device:deleted")
	DahuaEmailCreated  = system.NewEventBuilder[int64]("dahua-email:created")
)
