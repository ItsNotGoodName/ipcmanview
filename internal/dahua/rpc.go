package dahua

import (
	"context"
	"errors"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/license"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuarpc/modules/magicbox"
)

func RPCDetailGet(ctx context.Context, rpcClient dahuarpc.Client) (models.DahuaCameraDetail, error) {
	sn, err := magicbox.GetSerialNo(ctx, rpcClient)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	deviceClass, err := magicbox.GetDeviceClass(ctx, rpcClient)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	deviceType, err := magicbox.GetDeviceType(ctx, rpcClient)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	hardwareVersion, err := magicbox.GetHardwareVersion(ctx, rpcClient)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	marketArea, err := magicbox.GetMarketArea(ctx, rpcClient)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	ProcessInfo, err := magicbox.GetProcessInfo(ctx, rpcClient)
	if isNotResponseError(err) {
		return models.DahuaCameraDetail{}, err
	}

	vendor, err := magicbox.GetVendor(ctx, rpcClient)
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

func RPCSoftwareVersionGet(ctx context.Context, rpcClient dahuarpc.Client) (magicbox.GetSoftwareVersionResult, error) {
	res, err := magicbox.GetSoftwareVersion(ctx, rpcClient)
	if isNotResponseError(err) {
		return magicbox.GetSoftwareVersionResult{}, err
	}

	return res, nil
}

func RPCLicenseList(ctx context.Context, rpcClient dahuarpc.Client) ([]license.LicenseInfo, error) {
	res, err := license.GetLicenseInfo(ctx, rpcClient)
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
