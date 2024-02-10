package dahua

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
	sq "github.com/Masterminds/squirrel"
)

type dbCountRow struct {
	Count int64
}

// dbSelectFilter applies an authorization filter to a select query.
func dbSelectFilter(ctx context.Context, sb sq.SelectBuilder, deviceIDField string, levels ...models.DahuaPermissionLevel) sq.SelectBuilder {
	actor := core.UseActor(ctx)

	if actor.Admin {
		return sb
	}

	var level models.DahuaPermissionLevel
	if len(levels) != 0 {
		level = levels[0]
	}

	return sb.
		Where(sq.Expr(deviceIDField+` IN (
			SELECT
				device_id
			FROM
				dahua_permissions
			WHERE
				dahua_permissions.level > ?
				AND (
					dahua_permissions.user_id = ?
					OR dahua_permissions.group_id IN (
						SELECT
							group_id
						FROM
							group_users
						WHERE
							group_users.user_id = ?
					)
				)
			)
		`, level, actor.UserID, actor.UserID))
}

func GetConn(ctx context.Context, db sqlite.DB, id int64) (Conn, error) {
	actor := core.UseActor(ctx)
	v, err := db.C().DahuaGetConn(ctx, repo.DahuaGetConnParams{
		ID:     id,
		Admin:  actor.Admin,
		UserID: core.NewNullInt64(actor.UserID),
	})
	if err != nil {
		return Conn{}, err
	}

	return Conn{
		ID:       v.ID,
		URL:      v.Url.URL,
		Username: v.Username,
		Password: v.Password,
		Location: v.Location.Location,
		Feature:  v.Feature,
		Seed:     int(v.Seed),
	}, nil
}

func ListConn(ctx context.Context, db sqlite.DB) ([]Conn, error) {
	actor := core.UseActor(ctx)
	vv, err := db.C().DahuaListConn(ctx, repo.DahuaListConnParams{
		Admin:  actor.Admin,
		UserID: core.NewNullInt64(actor.UserID),
	})
	if err != nil {
		return nil, err
	}

	conns := make([]Conn, 0, len(vv))
	for _, v := range vv {
		conns = append(conns, Conn{
			ID:       v.ID,
			URL:      v.Url.URL,
			Username: v.Username,
			Password: v.Password,
			Location: v.Location.Location,
			Feature:  v.Feature,
			Seed:     int(v.Seed),
		})
	}

	return conns, nil
}

func ListDeviceIDs(ctx context.Context, db sqlite.DB) ([]int64, error) {
	sb := sq.
		Select("id").
		From("dahua_devices")

	var res []int64
	err := ssq.Query(ctx, db, &res, dbSelectFilter(ctx, sb, "dahua_devices.id"))
	return res, err
}

func CountFiles(ctx context.Context, db sqlite.DB) (int64, error) {
	sb := sq.
		Select("COUNT(*) AS count").
		From("dahua_files")

	var res dbCountRow
	err := ssq.QueryOne(ctx, db, &res, dbSelectFilter(ctx, sb, "dahua_files.device_id"))
	return res.Count, err
}

func CountEvents(ctx context.Context, db sqlite.DB) (int64, error) {
	sb := sq.
		Select("COUNT(*) AS count").
		From("dahua_events")

	var res dbCountRow
	err := ssq.QueryOne(ctx, db, &res, dbSelectFilter(ctx, sb, "dahua_events.device_id"))
	return res.Count, err
}

func CountEmails(ctx context.Context, db sqlite.DB) (int64, error) {
	sb := sq.
		Select("COUNT(*) AS count").
		From("dahua_email_messages")

	var res dbCountRow
	err := ssq.QueryOne(ctx, db, &res, dbSelectFilter(ctx, sb, "dahua_email_messages.device_id", models.DahuaPermissionLevelAdmin))
	return res.Count, err
}

type ListLatestEmailsResult struct {
	repo.DahuaEmailMessage
	AttachmentCount int64
}

func ListLatestEmails(ctx context.Context, db sqlite.DB, count int) ([]ListLatestEmailsResult, error) {
	sb := sq.
		Select("dahua_email_messages.*", "COUNT(dahua_email_attachments.id) AS attachment_count").
		From("dahua_email_messages").
		LeftJoin("dahua_email_attachments ON dahua_email_attachments.message_id = dahua_email_messages.id").
		OrderBy("created_at DESC").
		GroupBy("dahua_email_messages.id").
		Limit(uint64(count))

	var res []ListLatestEmailsResult
	err := ssq.Query(ctx, db, &res, dbSelectFilter(ctx, sb, "dahua_email_messages.device_id", models.DahuaPermissionLevelAdmin))
	return res, err
}

func ListLatestFiles(ctx context.Context, db sqlite.DB, count int) ([]repo.DahuaFile, error) {
	sb := sq.
		Select("*").
		From("dahua_files").
		OrderBy("start_time DESC").
		Limit(uint64(count))

	var res []repo.DahuaFile
	err := ssq.Query(ctx, db, &res, dbSelectFilter(ctx, sb, "dahua_files.device_id"))
	return res, err
}

type GetDeviceFilter struct {
	ID int64
	IP string
}

func GetDevice(ctx context.Context, db sqlite.DB, filter GetDeviceFilter) (repo.DahuaDevice, error) {
	eq := sq.Eq{}
	if filter.ID != 0 {
		eq["id"] = filter.ID
	}
	if filter.IP != "" {
		eq["ip"] = filter.IP
	}

	sb := sq.
		Select("*").
		From("dahua_devices").
		Where(eq)

	var res repo.DahuaDevice
	err := ssq.QueryOne(ctx, db, &res, dbSelectFilter(ctx, sb, "dahua_devices.id"))
	return res, err
}

func ListDevices(ctx context.Context, db sqlite.DB) ([]repo.DahuaDevice, error) {
	sb := sq.
		Select("*").
		From("dahua_devices")

	var res []repo.DahuaDevice
	err := ssq.Query(ctx, db, &res, dbSelectFilter(ctx, sb, "dahua_devices.id"))
	return res, err
}

type ListEmailsParams struct {
	pagination.Page
	Ascending bool
}

type ListEmailsResult struct {
	pagination.PageResult
	Items []ListEmailsResultItems
}

type ListEmailsResultItems struct {
	repo.DahuaEmailMessage
	DeviceName      string
	AttachmentCount int
}

func ListEmails(ctx context.Context, db sqlite.DB, arg ListEmailsParams) (ListEmailsResult, error) {
	order := "dahua_email_messages.id"
	if arg.Ascending {
		order += " ASC"
	} else {
		order += " DESC"
	}
	sb := sq.
		Select(
			"dahua_email_messages.*",
			"COUNT(dahua_email_attachments.id) AS attachment_count",
			"dahua_devices.name AS device_name",
		).
		From("dahua_email_messages").
		LeftJoin("dahua_email_attachments ON dahua_email_attachments.message_id = dahua_email_messages.id").
		LeftJoin("dahua_devices ON dahua_devices.id = dahua_email_messages.device_id").
		OrderBy(order).
		GroupBy("dahua_email_messages.id").
		Offset(uint64(arg.Offset())).
		Limit(uint64(arg.Limit()))

	var items []ListEmailsResultItems
	err := ssq.Query(ctx, db, &items, dbSelectFilter(ctx, sb, "dahua_email_messages.device_id", models.DahuaPermissionLevelAdmin))
	if err != nil {
		return ListEmailsResult{}, err
	}

	count, err := CountEmails(ctx, db)
	if err != nil {
		return ListEmailsResult{}, err
	}

	return ListEmailsResult{
		PageResult: arg.Result(int(count)),
		Items:      items,
	}, nil
}
