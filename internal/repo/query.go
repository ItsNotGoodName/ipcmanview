package repo

// import (
//
//	"context"
//	"database/sql"
//	"encoding/base64"
//	"errors"
//
//	"github.com/ItsNotGoodName/ipcmanview/internal/models"
//	"github.com/ItsNotGoodName/ipcmanview/internal/types"
//	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
//	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
//	sq "github.com/Masterminds/squirrel"
//
// )
//
//	type ListDahuaEventParams struct {
//		pagination.Page
//		Code      []string
//		Action    []string
//		DeviceID  []int64
//		Start     types.Time
//		End       types.Time
//		Ascending bool
//	}
//
//	type ListDahuaEventResult struct {
//		pagination.PageResult
//		Data []ListDahuaEvent
//	}
//
//	type ListDahuaEvent struct {
//		DeviceName string
//		DahuaEvent
//	}
//
//	func (db DB) ListDahuaEvent(ctx context.Context, arg ListDahuaEventParams) (ListDahuaEventResult, error) {
//		where := sq.And{}
//
//		eq := sq.Eq{}
//		if len(arg.Code) != 0 {
//			eq["code"] = arg.Code
//		}
//		if len(arg.Action) != 0 {
//			eq["action"] = arg.Action
//		}
//		if len(arg.DeviceID) != 0 {
//			eq["device_id"] = arg.DeviceID
//		}
//		where = append(where, eq)
//
//		and := sq.And{}
//		if !arg.Start.IsZero() {
//			and = append(and, sq.GtOrEq{"created_at": arg.Start})
//		}
//		if !arg.End.IsZero() {
//			and = append(and, sq.Lt{"created_at": arg.End})
//		}
//		where = append(where, and)
//
//		order := "created_at DESC"
//		if arg.Ascending {
//			order = "created_at ASC"
//		}
//
//		var res []ListDahuaEvent
//		err := ssq.Query(ctx, db, &res, sq.
//			Select("e.*, d.name as device_name").
//			From("dahua_events AS e").
//			Where(where).
//			OrderBy(order).
//			Limit(uint64(arg.Page.Limit())).
//			Offset(uint64(arg.Page.Offset())).
//			LeftJoin("dahua_devices AS d ON d.id = e.device_id"))
//		if err != nil {
//			return ListDahuaEventResult{}, err
//		}
//
//		var count int
//		err = ssq.QueryOne(ctx, db, &count, sq.
//			Select("COUNT(*)").
//			From("dahua_events").
//			Where(where))
//		if err != nil {
//			return ListDahuaEventResult{}, err
//		}
//
//		return ListDahuaEventResult{
//			PageResult: arg.Page.Result(count),
//			Data:       res,
//		}, nil
//	}
//
//	type DahuaFileFilter struct {
//		Type      []string
//		DeviceID  []int64
//		Start     types.Time
//		End       types.Time
//		Ascending bool
//		Storage   []models.Storage
//	}
//
//	func (arg DahuaFileFilter) where(where sq.And) sq.And {
//		eq := sq.Eq{}
//		if len(arg.Type) != 0 {
//			eq["type"] = arg.Type
//		}
//		if len(arg.DeviceID) != 0 {
//			eq["device_id"] = arg.DeviceID
//		}
//		if len(arg.Storage) != 0 {
//			eq["storage"] = arg.Storage
//		}
//		where = append(where, eq)
//
//		and := sq.And{}
//		if !arg.Start.IsZero() {
//			and = append(and, sq.GtOrEq{"start_time": arg.Start})
//		}
//		if !arg.End.IsZero() {
//			and = append(and, sq.Lt{"start_time": arg.End})
//		}
//		where = append(where, and)
//
//		return where
//	}
//
//	func (arg DahuaFileFilter) order() string {
//		if arg.Ascending {
//			return "start_time ASC"
//		} else {
//			return "start_time DESC"
//		}
//	}
//
//	type ListDahuaFileParams struct {
//		pagination.Page
//		DahuaFileFilter
//	}
//
//	type ListDahuaFileResult struct {
//		pagination.PageResult
//		Data []DahuaFile
//	}
//
//	func (db DB) ListDahuaFile(ctx context.Context, arg ListDahuaFileParams) (ListDahuaFileResult, error) {
//		where := arg.where(sq.And{})
//		order := arg.order()
//
//		var res []DahuaFile
//		err := ssq.Query(ctx, db, &res, sq.
//			Select("*").
//			From("dahua_files").
//			Where(where).
//			OrderBy(order).
//			Limit(uint64(arg.Page.Limit())).
//			Offset(uint64(arg.Page.Offset())))
//		if err != nil {
//			return ListDahuaFileResult{}, err
//		}
//
//		var count int
//		err = ssq.QueryOne(ctx, db, &count, sq.
//			Select("COUNT(*)").
//			From("dahua_files").
//			Where(where))
//		if err != nil {
//			return ListDahuaFileResult{}, err
//		}
//
//		return ListDahuaFileResult{
//			PageResult: arg.Page.Result(count),
//			Data:       res,
//		}, nil
//	}
//
//	type CursorListDahuaFileParams struct {
//		Cursor  string
//		PerPage int
//		DahuaFileFilter
//	}
//
//	type CursorListDahuaFileResult struct {
//		Cursor  string
//		HasMore bool
//		Data    []DahuaFile
//	}
//
//	func (db DB) CursorListDahuaFile(ctx context.Context, arg CursorListDahuaFileParams) (CursorListDahuaFileResult, error) {
//		where := arg.where(sq.And{})
//		if arg.Cursor != "" {
//			b, err := base64.URLEncoding.DecodeString(arg.Cursor)
//			if err != nil {
//				return CursorListDahuaFileResult{}, err
//			}
//
//			var startTime types.Time
//
//			err = startTime.UnmarshalBinary(b)
//			if err != nil {
//				return CursorListDahuaFileResult{}, err
//			}
//
//			if arg.Ascending {
//				where = append(where, sq.GtOrEq{"start_time": startTime})
//			} else {
//				where = append(where, sq.LtOrEq{"start_time": startTime})
//			}
//		}
//
//		order := arg.order()
//		limit := arg.PerPage + 1
//
//		var res []DahuaFile
//		err := ssq.Query(ctx, db, &res, sq.
//			Select("*").
//			From("dahua_files").
//			Where(where).
//			OrderBy(order).
//			Limit(uint64(limit)))
//		if err != nil {
//			return CursorListDahuaFileResult{}, err
//		}
//		length := len(res)
//
//		if length == 0 || length != limit {
//			return CursorListDahuaFileResult{
//				Cursor:  "",
//				HasMore: false,
//				Data:    res,
//			}, nil
//		}
//
//		data, last := res[:length-1], res[length-1]
//
//		var cursor string
//		{
//			b, err := last.StartTime.MarshalBinary()
//			if err != nil {
//				return CursorListDahuaFileResult{}, nil
//			}
//
//			cursor = base64.URLEncoding.EncodeToString(b)
//		}
//
//		return CursorListDahuaFileResult{
//			Cursor:  cursor,
//			HasMore: true,
//			Data:    data,
//		}, nil
//	}
//
//
// type ListDahuaDeviceByFeatureRow = listDahuaDeviceByFeatureRow
//
//	func (db DB) ListDahuaDeviceByFeature(ctx context.Context, features ...models.DahuaFeature) ([]ListDahuaDeviceByFeatureRow, error) {
//		var feature models.DahuaFeature
//		for _, f := range features {
//			feature = feature | f
//		}
//		return db.listDahuaDeviceByFeature(ctx, feature)
//	}
//
