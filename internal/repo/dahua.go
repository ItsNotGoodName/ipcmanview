package repo

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
	sq "github.com/Masterminds/squirrel"
)

func dahuaSelectFilter(ctx context.Context, sb sq.SelectBuilder, joinTableField string) sq.SelectBuilder {
	actor := core.UseActor(ctx)

	if actor.Admin {
		return sb
	}

	return sb.
		Join("dahua_permissions ON dahua_permissions.device_id = " + joinTableField).
		Where(sq.Expr(`
			dahua_permissions.user_id = ?
			OR dahua_permissions.group_id IN (
				SELECT
					group_id
				FROM
					group_users
				WHERE
					group_users.user_id = ?
			)
		`, actor.UserID, actor.UserID))
}

// DahuaFatDevice

type DahuaFatDevice struct {
	DahuaDevice
	Seed int64
}

type DahuaFatDeviceParams struct {
	IPs      []string
	IDs      []int64
	Features []models.DahuaFeature
	Limit    int
}

func (db DB) DahuaListFatDevices(ctx context.Context, args ...DahuaFatDeviceParams) ([]DahuaFatDevice, error) {
	var arg DahuaFatDeviceParams
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

	sb = sb.Where(and)
	// LIMIT
	if arg.Limit != 0 {
		sb = sb.Limit(uint64(arg.Limit))
	}

	var res []DahuaFatDevice
	return res, ssq.Query(ctx, db, &res, dahuaSelectFilter(ctx, sb, "dahua_devices.id"))
}

func (db DB) DahuaGetFatDevice(ctx context.Context, arg DahuaFatDeviceParams) (DahuaFatDevice, error) {
	arg.Limit = 1
	devices, err := db.DahuaListFatDevices(ctx, arg)
	if err != nil {
		return DahuaFatDevice{}, err
	}
	if len(devices) == 0 {
		return DahuaFatDevice{}, ErrNotFound
	}
	return devices[0], nil
}

// DahuaDevicePermission

// var dahuaListDahuaDevicePermissions = fmt.Sprintf(`
// SELECT
//   d.id as device_id,
//   CASE
//     WHEN a.user_id IS NOT NULL THEN %d
//     ELSE max(p.level)
//   END AS level
// FROM
//   dahua_devices AS d
//   LEFT JOIN dahua_permissions AS p ON p.device_id = d.id
//   LEFT JOIN admins AS a ON a.user_id = ?1
// WHERE
//   -- Allow if user is admin
//   a.user_id IS NOT NULL
//   -- Filter by level
//   OR (
//     -- Allow if user owns the permission
//     p.user_id = ?1
//     -- Allow if user is a part of the group that owns the permission
//     OR p.group_id IN (
//       SELECT
//         group_id
//       FROM
//         group_users
//       WHERE
//         group_users.user_id = ?1
//     )
//   )
// GROUP BY
//   d.id
// 	`, models.DahuaPermissionLevelAdmin)
//
// func (db DB) DahuaListDahuaDevicePermissions(ctx context.Context, userID int64) (models.DahuaDevicePermissions, error) {
// 	rows, err := db.QueryContext(ctx, dahuaListDahuaDevicePermissions, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	var res []models.DahuaDevicePermission
// 	if err := sqlscan.ScanAll(&res, rows); err != nil {
// 		return nil, err
// 	}
//
// 	return res, nil
// }
