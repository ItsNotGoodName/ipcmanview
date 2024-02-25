package event

import (
	"encoding/json"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

const (
	ActionDahuaDeviceCreated  models.EventAction = "dahua-device:created"
	ActionDahuaDeviceUpdated  models.EventAction = "dahua-device:updated"
	ActionDahuaDeviceDeleted  models.EventAction = "dahua-device:deleted"
	ActionDahuaEmailCreated   models.EventAction = "dahua-email:created"
	ActionUserSecurityUpdated models.EventAction = "user-security:updated"
)

func DataAsInt64(evt repo.Event) int64 {
	var deviceID int64
	err := json.Unmarshal(evt.Data.RawMessage, &deviceID)
	if err != nil {
		return 0
	}
	return deviceID
}
