package dahua

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
)

func ListIDsByLevel(rows []repo.DahuaListDevicePermissionsRow, level models.DahuaPermissionLevel) []int64 {
	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		if row.Level >= level {
			ids = append(ids, row.ID)
		}
	}
	return ids
}
