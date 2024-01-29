package repo

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
	sq "github.com/Masterminds/squirrel"
)

type FatDahuaDeviceParams struct {
	IPs      []string
	IDs      []int64
	Features []models.DahuaFeature
}

type FatDahuaDevice struct {
	DahuaDevice
	Seed int64
}

func (db DB) DahuaListDevices(ctx context.Context, args ...FatDahuaDeviceParams) ([]FatDahuaDevice, error) {
	var arg FatDahuaDeviceParams
	if len(args) != 0 {
		arg = args[0]
	}

	// SELECT ...
	sb := sq.
		Select(
			"dahua_devices.*",
			"coalesce(seed, id) AS seed",
		).
		From("dahua_devices").
		LeftJoin("dahua_seeds ON dahua_seeds.device_id = dahua_devices.id")
	// WHERE
	and := sq.And{}

	eq := sq.Eq{}
	if arg.IPs != nil {
		eq["ip"] = arg.IPs
	}
	if arg.IDs != nil {
		eq["id"] = arg.IDs
	}
	and = append(and, eq)

	if len(arg.Features) != 0 {
		var feature models.DahuaFeature
		for _, f := range arg.Features {
			feature = feature | f
		}
		and = append(and, sq.Expr("feature & ? = ?", feature, feature))
	}

	sb = sb.Where(eq)

	var res []FatDahuaDevice
	return res, ssq.Query(ctx, db, &res, sb)
}

func (db DB) DahuaGetDevice(ctx context.Context, arg FatDahuaDeviceParams) (FatDahuaDevice, error) {
	devices, err := db.DahuaListDevices(ctx, arg)
	if err != nil {
		return FatDahuaDevice{}, err
	}
	if len(devices) == 0 {
		return FatDahuaDevice{}, ErrNotFound
	}
	return devices[0], nil
}
