package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmango/internal/models"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/modules/license"
	"github.com/ItsNotGoodName/ipcmango/pkg/qes"
)

type WorkerJob interface {
	Execute(ctx context.Context, gen dahua.GenRPC, cam models.DahuaCamera) error
}

var _ WorkerJob = (*TestJob)(nil)

type TestJob struct {
	DB qes.Querier
}

func (t *TestJob) Execute(ctx context.Context, gen dahua.GenRPC, cam models.DahuaCamera) error {
	details, err := CameraDetailGet(ctx, gen)
	if err != nil {
		return err
	}

	if err := DB.CameraDetailUpdate(ctx, t.DB, cam.ID, details); err != nil {
		return err
	}

	softwares, err := CameraSoftwareVersionGet(ctx, gen)
	if err != nil {
		return err
	}

	if err := DB.CameraSoftwaresUpdate(ctx, t.DB, cam.ID, softwares); err != nil {
		return err
	}

	licenses, err := license.GetLicenseInfo(ctx, gen)
	if err != nil {
		return err
	}

	if err := DB.CameraLicensesReplace(ctx, t.DB, cam.ID, licenses); err != nil {
		return err
	}

	return nil
}
