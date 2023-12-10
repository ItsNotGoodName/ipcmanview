package webserver

import (
	"cmp"
	"context"
	"slices"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/rs/zerolog/log"
)

func useDahuaTables(ctx context.Context, db repo.DB, dahuaStore *dahua.Store) (any, error) {
	dbCameras, err := db.ListDahuaCamera(ctx)
	if err != nil {
		return nil, err
	}

	type cameraData struct {
		detail               models.DahuaDetail
		softwareVersion      models.DahuaSoftwareVersion
		licenses             []models.DahuaLicense
		storage              []models.DahuaStorage
		coaxialcontrolStatus []models.DahuaCoaxialStatus
	}

	cameras := make([]models.DahuaCamera, 0, len(dbCameras))
	for _, row := range dbCameras {
		cameras = append(cameras, row.Convert())
	}
	conns := dahuaStore.ConnList(ctx, cameras)

	cameraDataC := make(chan cameraData, len(dbCameras))
	wg := sync.WaitGroup{}
	for _, conn := range conns {
		wg.Add(1)
		go func(conn dahua.Conn) {
			defer wg.Done()

			log := log.With().Int64("id", conn.Camera.ID).Logger()

			var data cameraData

			{
				res, err := dahua.GetDahuaDetail(ctx, conn.Camera.ID, conn.RPC)
				if err != nil {
					log.Err(err).Msg("Failed to get detail")
					return
				}

				data.detail = res
			}

			{
				res, err := dahua.GetSoftwareVersion(ctx, conn.Camera.ID, conn.RPC)
				if err != nil {
					log.Err(err).Msg("Failed to get software version")
					return
				}

				data.softwareVersion = res
			}

			{
				res, err := dahua.GetLicenseList(ctx, conn.Camera.ID, conn.RPC)
				if err != nil {
					log.Err(err).Msg("Failed to get licenses")
					return
				}

				data.licenses = res
			}

			{
				res, err := dahua.GetStorage(ctx, conn.Camera.ID, conn.RPC)
				if err != nil {
					log.Err(err).Msg("Failed to get storage")
				}

				data.storage = res
			}

			{
				caps, err := dahua.GetCoaxialCaps(ctx, conn.Camera.ID, conn.RPC, 0)
				if err != nil {
					log.Err(err).Msg("Failed to get coaxial caps")
					return
				}

				if caps.SupportControlLight || caps.SupportControlSpeaker || caps.SupportControlFullcolorLight {
					res, err := dahua.GetCoaxialStatus(ctx, conn.Camera.ID, conn.RPC, 0)
					if err != nil {
						log.Err(err).Msg("Failed to get coaxial status")
						return
					}

					data.coaxialcontrolStatus = append(data.coaxialcontrolStatus, res)
				}
			}

			cameraDataC <- data
		}(conn)
	}
	wg.Wait()
	close(cameraDataC)

	status := make([]models.DahuaStatus, 0, len(conns))
	for _, conn := range conns {
		status = append(status, dahua.GetDahuaStatus(conn.Camera, conn.RPC))
	}

	details := make([]models.DahuaDetail, 0, len(dbCameras))
	softwareVersions := make([]models.DahuaSoftwareVersion, 0, len(dbCameras))
	licenses := make([]models.DahuaLicense, 0, len(dbCameras))
	storage := make([]models.DahuaStorage, 0, len(dbCameras))
	coaxialStatus := make([]models.DahuaCoaxialStatus, 0, len(dbCameras))
	for data := range cameraDataC {
		if data.detail.CameraID != 0 {
			details = append(details, data.detail)
		}
		if data.softwareVersion.CameraID != 0 {
			softwareVersions = append(softwareVersions, data.softwareVersion)
		}
		licenses = append(licenses, data.licenses...)
		storage = append(storage, data.storage...)
		coaxialStatus = append(coaxialStatus, data.coaxialcontrolStatus...)
	}
	slices.SortFunc(details, func(a, b models.DahuaDetail) int { return cmp.Compare(a.CameraID, b.CameraID) })
	slices.SortFunc(softwareVersions, func(a, b models.DahuaSoftwareVersion) int { return cmp.Compare(a.CameraID, b.CameraID) })
	slices.SortFunc(licenses, func(a, b models.DahuaLicense) int { return cmp.Compare(a.CameraID, b.CameraID) })
	slices.SortFunc(storage, func(a, b models.DahuaStorage) int { return cmp.Compare(a.CameraID, b.CameraID) })
	slices.SortFunc(coaxialStatus, func(a, b models.DahuaCoaxialStatus) int { return cmp.Compare(a.CameraID, b.CameraID) })

	return Data{
		"Cameras":          dbCameras,
		"Status":           status,
		"Details":          details,
		"SoftwareVersions": softwareVersions,
		"Licenses":         licenses,
		"Storage":          storage,
		"CoaxialStatus":    coaxialStatus,
	}, nil
}
