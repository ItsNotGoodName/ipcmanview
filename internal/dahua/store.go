package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
)

type Store interface {
	ClientRPC(ctx context.Context, cameraID int64) (dahuarpc.Client, error)
	// ClientCGI(ctx context.Context, cameraID int64) (dahuacgi.Client, error)
}
