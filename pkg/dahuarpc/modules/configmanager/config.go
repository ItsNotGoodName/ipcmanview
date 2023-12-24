package configmanager

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/configmanager/config"
)

func VideoInMode(ctx context.Context, c dahuarpc.Conn) (Config[config.VideoInMode], error) {
	return GetConfig[config.VideoInMode](ctx, c, "VideoInMode", true)
}
