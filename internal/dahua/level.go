package dahua

import (
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

const levelDefault = models.DahuaPermissionLevelUser
const levelEmail = models.DahuaPermissionLevelUser

// type Action string
//
// var (
// 	ActionDetailRead Action = "details.read"
// 	ActionEmailsRead Action = "emails.read"
// )
//
// func (action Action) Can(ctx context.Context, level models.DahuaPermissionLevel) bool {
// 	if level == models.DahuaPermissionLevelAdmin {
// 		return true
// 	}
//
// 	session, ok := auth.UseSession(ctx)
// 	if ok && session.Admin {
// 		return true
// 	}
//
// 	if strings.HasSuffix(string(action), ".read") {
// 		return true
// 	}
//
// 	return false
// }
