package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
	sq "github.com/Masterminds/squirrel"
)

type CreateDahuaCameraParams = createDahuaCameraParams

type CreateDahuaFileCursorParams = createDahuaFileCursorParams

func (db DB) CreateDahuaCamera(ctx context.Context, arg CreateDahuaCameraParams, args2 CreateDahuaFileCursorParams) (int64, error) {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return 0, nil
	}
	defer tx.Rollback()

	id, err := tx.createDahuaCamera(ctx, arg)
	if err != nil {
		return 0, err
	}

	// TODO: sql.NullInt64 should just be int64...
	err = tx.allocateDahuaSeed(ctx, sql.NullInt64{
		Valid: true,
		Int64: id,
	})
	if err != nil {
		return 0, err
	}

	args2.CameraID = id
	err = tx.createDahuaFileCursor(ctx, args2)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db DB) UpsertDahuaFiles(ctx context.Context, args CreateDahuaFileParams) (int64, error) {
	id, err := db.UpdateDahuaFile(ctx, UpdateDahuaFileParams{
		CameraID:    args.CameraID,
		Channel:     args.Channel,
		StartTime:   args.StartTime,
		EndTime:     args.EndTime,
		Length:      args.Length,
		Type:        args.Type,
		FilePath:    args.FilePath,
		Duration:    args.Duration,
		Disk:        args.Disk,
		VideoStream: args.VideoStream,
		Flags:       args.Flags,
		Events:      args.Events,
		Cluster:     args.Cluster,
		Partition:   args.Partition,
		PicIndex:    args.PicIndex,
		Repeat:      args.Repeat,
		WorkDir:     args.WorkDir,
		WorkDirSn:   args.WorkDirSn,
		UpdatedAt:   args.UpdatedAt,
	})
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	return db.CreateDahuaFile(ctx, args)
}

type ListDahuaEventParams struct {
	pagination.Page
	Code      []string
	Action    []string
	CameraID  []int64
	Start     types.Time
	End       types.Time
	Ascending bool
}

type ListDahuaEventResult struct {
	pagination.PageResult
	Data []DahuaEvent
}

func (db DB) ListDahuaEvent(ctx context.Context, arg ListDahuaEventParams) (ListDahuaEventResult, error) {
	where := sq.And{}

	eq := sq.Eq{}
	if len(arg.Code) != 0 {
		eq["code"] = arg.Code
	}
	if len(arg.Action) != 0 {
		eq["action"] = arg.Action
	}
	if len(arg.CameraID) != 0 {
		eq["camera_id"] = arg.CameraID
	}
	where = append(where, eq)

	and := sq.And{}
	if !arg.Start.IsZero() {
		and = append(and, sq.GtOrEq{"created_at": arg.Start})
	}
	if !arg.End.IsZero() {
		and = append(and, sq.Lt{"created_at": arg.End})
	}
	where = append(where, and)

	order := "created_at DESC"
	if arg.Ascending {
		order = "created_at ASC"
	}

	var res []DahuaEvent
	err := ssq.Query(ctx, db, &res, sq.
		Select("*").
		From("dahua_events").
		Where(where).
		OrderBy(order).
		Limit(uint64(arg.Page.Limit())).
		Offset(uint64(arg.Page.Offset())))
	if err != nil {
		return ListDahuaEventResult{}, err
	}

	var count int
	err = ssq.QueryOne(ctx, db, &count, sq.
		Select("COUNT(*)").
		From("dahua_events").
		Where(where))
	if err != nil {
		return ListDahuaEventResult{}, err
	}

	return ListDahuaEventResult{
		PageResult: arg.Page.Result(count),
		Data:       res,
	}, nil
}
