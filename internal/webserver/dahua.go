package webserver

import (
	"cmp"
	"context"
	"slices"
	"sync"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahuacore"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/rs/zerolog/log"
)

func useDahuaTables(ctx context.Context, db repo.DB, dahuaStore *dahuacore.Store) (any, error) {
	dbDevices, err := db.ListDahuaDevice(ctx)
	if err != nil {
		return nil, err
	}

	type deviceData struct {
		detail               models.DahuaDetail
		softwareVersion      models.DahuaSoftwareVersion
		licenses             []models.DahuaLicense
		storage              []models.DahuaStorage
		coaxialcontrolStatus []models.DahuaCoaxialStatus
	}

	devices := make([]models.DahuaConn, 0, len(dbDevices))
	for _, row := range dbDevices {
		devices = append(devices, row.Convert().DahuaConn)
	}
	conns := dahuaStore.ConnList(ctx, devices)

	deviceDataC := make(chan deviceData, len(dbDevices))
	wg := sync.WaitGroup{}
	for _, conn := range conns {
		wg.Add(1)
		go func(conn dahuacore.Conn) {
			defer wg.Done()

			log := log.With().Int64("id", conn.Device.ID).Logger()

			var data deviceData

			{
				res, err := dahuacore.GetDahuaDetail(ctx, conn.Device.ID, conn.RPC)
				if err != nil {
					log.Err(err).Msg("Failed to get detail")
					return
				}

				data.detail = res
			}

			{
				res, err := dahuacore.GetSoftwareVersion(ctx, conn.Device.ID, conn.RPC)
				if err != nil {
					log.Err(err).Msg("Failed to get software version")
					return
				}

				data.softwareVersion = res
			}

			{
				res, err := dahuacore.GetLicenseList(ctx, conn.Device.ID, conn.RPC)
				if err != nil {
					log.Err(err).Msg("Failed to get licenses")
					return
				}

				data.licenses = res
			}

			{
				res, err := dahuacore.GetStorage(ctx, conn.Device.ID, conn.RPC)
				if err != nil {
					log.Err(err).Msg("Failed to get storage")
				}

				data.storage = res
			}

			{
				caps, err := dahuacore.GetCoaxialCaps(ctx, conn.Device.ID, conn.RPC, 1)
				if err != nil {
					log.Err(err).Msg("Failed to get coaxial caps")
					return
				}

				if caps.SupportControlLight || caps.SupportControlSpeaker || caps.SupportControlFullcolorLight {
					res, err := dahuacore.GetCoaxialStatus(ctx, conn.Device.ID, conn.RPC, 1)
					if err != nil {
						log.Err(err).Msg("Failed to get coaxial status")
						return
					}

					data.coaxialcontrolStatus = append(data.coaxialcontrolStatus, res)
				}
			}

			deviceDataC <- data
		}(conn)
	}
	wg.Wait()
	close(deviceDataC)

	status := make([]models.DahuaStatus, 0, len(conns))
	for _, conn := range conns {
		status = append(status, dahuacore.GetDahuaStatus(conn.Device, conn.RPC))
	}

	details := make([]models.DahuaDetail, 0, len(dbDevices))
	softwareVersions := make([]models.DahuaSoftwareVersion, 0, len(dbDevices))
	licenses := make([]models.DahuaLicense, 0, len(dbDevices))
	storage := make([]models.DahuaStorage, 0, len(dbDevices))
	coaxialStatus := make([]models.DahuaCoaxialStatus, 0, len(dbDevices))
	for data := range deviceDataC {
		if data.detail.DeviceID != 0 {
			details = append(details, data.detail)
		}
		if data.softwareVersion.DeviceID != 0 {
			softwareVersions = append(softwareVersions, data.softwareVersion)
		}
		licenses = append(licenses, data.licenses...)
		storage = append(storage, data.storage...)
		coaxialStatus = append(coaxialStatus, data.coaxialcontrolStatus...)
	}
	slices.SortFunc(details, func(a, b models.DahuaDetail) int { return cmp.Compare(a.DeviceID, b.DeviceID) })
	slices.SortFunc(softwareVersions, func(a, b models.DahuaSoftwareVersion) int { return cmp.Compare(a.DeviceID, b.DeviceID) })
	slices.SortFunc(licenses, func(a, b models.DahuaLicense) int { return cmp.Compare(a.DeviceID, b.DeviceID) })
	slices.SortFunc(storage, func(a, b models.DahuaStorage) int { return cmp.Compare(a.DeviceID, b.DeviceID) })
	slices.SortFunc(coaxialStatus, func(a, b models.DahuaCoaxialStatus) int { return cmp.Compare(a.DeviceID, b.DeviceID) })

	return Data{
		"Devices":          dbDevices,
		"Status":           status,
		"Details":          details,
		"SoftwareVersions": softwareVersions,
		"Licenses":         licenses,
		"Storage":          storage,
		"CoaxialStatus":    coaxialStatus,
	}, nil
}
