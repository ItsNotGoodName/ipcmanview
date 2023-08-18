package rpc

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	"github.com/ItsNotGoodName/ipcmanview/server/rpcgen"
)

type DahuaService struct {
	db        qes.Querier
	super     *dahua.Supervisor
	scanSuper *dahua.ScanSupervisor
}

func NewDahuaService(db qes.Querier, super *dahua.Supervisor, scanSuper *dahua.ScanSupervisor) DahuaService {
	return DahuaService{
		db:        db,
		super:     super,
		scanSuper: scanSuper,
	}
}

var _ rpcgen.DahuaService = (*DahuaService)(nil)

// ActiveScannerCount implements rpcgen.DahuaService.
func (s DahuaService) ActiveScannerCount(ctx context.Context) (int, int, error) {
	count, err := dahua.DB.ScanActiveTaskCount(ctx, s.db)
	if err != nil {
		return 0, 0, err
	}

	return count, s.scanSuper.WorkerCount(), nil
}

// CameraCount implements rpcgen.DahuaService.
func (ds DahuaService) CameraCount(ctx context.Context) (int, error) {
	count, err := dahua.DB.CameraCount(ctx, ds.db)
	if err != nil {
		return 0, handleErr(rpcgen.ErrWebrpcInternalError, err)
	}

	return count, nil
}
