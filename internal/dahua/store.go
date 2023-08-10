package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahua"
)

type Store interface {
	GetGenRPC(ctx context.Context, cameraID int64) (dahua.GenRPC, error)
}
