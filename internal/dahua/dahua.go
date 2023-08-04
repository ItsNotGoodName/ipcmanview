package dahua

import (
	"context"
	"errors"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/modules/license"
	"github.com/ItsNotGoodName/ipcmango/pkg/dahua/modules/magicbox"
)

func CameraDetailGet(ctx context.Context, actor ActorHandle) (core.DahuaCameraDetail, error) {
	sn, err := magicbox.GetSerialNo(ctx, actor)
	if isNotResponseError(err) {
		return core.DahuaCameraDetail{}, err
	}

	deviceClass, err := magicbox.GetDeviceClass(ctx, actor)
	if isNotResponseError(err) {
		return core.DahuaCameraDetail{}, err
	}

	deviceType, err := magicbox.GetDeviceType(ctx, actor)
	if isNotResponseError(err) {
		return core.DahuaCameraDetail{}, err
	}

	hardwareVersion, err := magicbox.GetHardwareVersion(ctx, actor)
	if isNotResponseError(err) {
		return core.DahuaCameraDetail{}, err
	}

	marketArea, err := magicbox.GetMarketArea(ctx, actor)
	if isNotResponseError(err) {
		return core.DahuaCameraDetail{}, err
	}

	ProcessInfo, err := magicbox.GetProcessInfo(ctx, actor)
	if isNotResponseError(err) {
		return core.DahuaCameraDetail{}, err
	}

	vendor, err := magicbox.GetVendor(ctx, actor)
	if isNotResponseError(err) {
		return core.DahuaCameraDetail{}, err
	}

	return core.DahuaCameraDetail{
		ID:              actor.cam.ID,
		SN:              sn,
		DeviceClass:     deviceClass,
		DeviceType:      deviceType,
		HardwareVersion: hardwareVersion,
		MarketArea:      marketArea,
		ProcessInfo:     ProcessInfo,
		Vendor:          vendor,
	}, nil
}

func CameraSoftwareVersionGet(ctx context.Context, actor ActorHandle) (core.DahuaSoftwareVersion, error) {
	res, err := magicbox.GetSoftwareVersion(ctx, actor)
	if isNotResponseError(err) {
		return core.DahuaSoftwareVersion{}, err
	}

	return core.DahuaSoftwareVersion{
		ID:                      actor.cam.ID,
		Build:                   res.Build,
		BuildDate:               res.BuildDate,
		SecurityBaseLineVersion: res.SecurityBaseLineVersion,
		Version:                 res.Version,
		WebVersion:              res.WebVersion,
	}, nil
}

func LicensesList(ctx context.Context, actor ActorHandle) ([]core.DahuaCameraLicense, error) {
	res, err := license.GetLicenseInfo(ctx, actor)
	if isNotResponseError(err) {
		return []core.DahuaCameraLicense{}, err
	}

	licenses := make([]core.DahuaCameraLicense, 0, len(res))
	for _, v := range res {
		licenses = append(licenses, core.DahuaCameraLicense{
			ID:            actor.cam.ID,
			AbroadInfo:    v.AbroadInfo,
			AllType:       v.AllType,
			DigitChannel:  v.DigitChannel,
			EffectiveDays: v.EffectiveDays,
			EffectiveTime: v.EffectiveTime,
			LicenseID:     v.LicenseID,
			ProductType:   v.ProductType,
			Status:        v.Status,
			Username:      v.Username,
		})
	}

	return licenses, nil
}

func isNotResponseError(err error) bool {
	if err == nil {
		return false
	}
	var responseErr *dahua.ErrResponse
	return !errors.As(err, &responseErr)
}
