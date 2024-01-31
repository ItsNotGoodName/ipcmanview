package dahua

import (
	"context"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

type Action string

var (
	ActionDetailsRead Action = "details.read"
)

func (action Action) Can(ctx context.Context, level models.DahuaPermissionLevel) bool {
	if level == models.DahuaPermissionLevelAdmin {
		return true
	}

	session, ok := auth.UseSession(ctx)
	if ok && session.Admin {
		return true
	}

	if strings.HasSuffix(string(action), ".read") {
		return true
	}

	return false
}
