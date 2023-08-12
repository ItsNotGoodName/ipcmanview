package dahua

import (
	"context"
	"errors"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/license"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/magicbox"
)

func RPCDetailGet(ctx context.Context, gen dahuarpc.Gen) (models.DahuaCameraDetail, error) {
	sn, err := magicbox.GetSerialNo(ctx, gen)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	deviceClass, err := magicbox.GetDeviceClass(ctx, gen)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	deviceType, err := magicbox.GetDeviceType(ctx, gen)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	hardwareVersion, err := magicbox.GetHardwareVersion(ctx, gen)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	marketArea, err := magicbox.GetMarketArea(ctx, gen)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	ProcessInfo, err := magicbox.GetProcessInfo(ctx, gen)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	vendor, err := magicbox.GetVendor(ctx, gen)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	return models.DahuaCameraDetail{
		SN:              sn,
		DeviceClass:     deviceClass,
		DeviceType:      deviceType,
		HardwareVersion: hardwareVersion,
		MarketArea:      marketArea,
		ProcessInfo:     ProcessInfo,
		Vendor:          vendor,
	}, nil
}

func RPCSoftwareVersionGet(ctx context.Context, gen dahuarpc.Gen) (magicbox.GetSoftwareVersionResult, error) {
	res, err := magicbox.GetSoftwareVersion(ctx, gen)
	if isNotResponseError(err) {
		return magicbox.GetSoftwareVersionResult{}, err
	}

	return res, nil
}

func RPCLicenseList(ctx context.Context, gen dahuarpc.Gen) ([]license.LicenseInfo, error) {
	res, err := license.GetLicenseInfo(ctx, gen)
	if isNotResponseError(err) {
		return nil, err
	}

	return res, nil
}

func isNotResponseError(err error) bool {
	if err == nil {
		return false
	}
	var responseErr *dahuarpc.ErrResponse
	return !errors.As(err, &responseErr)
}
