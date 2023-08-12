package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type Store interface {
	GetGenRPC(ctx context.Context, cameraID int64) (dahuarpc.Gen, error)
}
