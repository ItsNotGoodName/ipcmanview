package dahua

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dbgen/postgres/dahua/model"
	dahua "github.com/ItsNotGoodName/ipcmanview/internal/dbgen/postgres/dahua/table"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahua/modules/license"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahua/modules/magicbox"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahua/modules/mediafilefind"
	"github.com/ItsNotGoodName/ipcmanview/pkg/qes"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5"
)

type dbT struct{}

var DB dbT

var dbCameraProjection ProjectionList = []Projection{
	dahua.Cameras.ID.AS("id"),
	dahua.Cameras.Address.AS("address"),
	dahua.Cameras.Username.AS("username"),
	dahua.Cameras.Password.AS("password"),
	dahua.Cameras.Location.AS("location"),
	dahua.Cameras.CreatedAt.AS("created_at"),
}

// CameraCreate [final]
func (dbT) CameraCreate(ctx context.Context, db qes.Querier, r models.DahuaCamera) (models.DahuaCamera, error) {
	if err := validateCamera(r); err != nil {
		return models.DahuaCamera{}, err
	}

	var camera models.DahuaCamera
	err := qes.ScanOne(ctx, db, &camera, dahua.Cameras.
		INSERT(
			dahua.Cameras.Address,
			dahua.Cameras.Username,
			dahua.Cameras.Password,
			dahua.Cameras.Location,
		).
		MODEL(struct {
			model.Cameras
			Location models.Location
		}{
			Cameras: model.Cameras{
				Address:  r.Address,
				Username: r.Username,
				Password: r.Password,
			},
			Location: r.Location,
		}).
		RETURNING(dbCameraProjection))
	return camera, err
}

// CameraUpdate [final]
func (dbT) CameraUpdate(ctx context.Context, db qes.Querier, r models.DahuaCamera) (models.DahuaCamera, error) {
	if err := validateCamera(r); err != nil {
		return models.DahuaCamera{}, err
	}

	var camera models.DahuaCamera
	err := qes.ScanOne(ctx, db, &camera, dahua.Cameras.
		UPDATE(
			dahua.Cameras.Address,
			dahua.Cameras.Username,
			dahua.Cameras.Password,
			dahua.Cameras.Location,
		).
		MODEL(struct {
			model.Cameras
			Location models.Location
		}{
			Cameras: model.Cameras{
				Address:  r.Address,
				Username: r.Username,
				Password: r.Password,
			},
			Location: r.Location,
		}).
		WHERE(dahua.Cameras.ID.EQ(Int64(r.ID))).
		RETURNING(dbCameraProjection))
	return camera, err
}

// CameraExists [final]
func (dbT) CameraExists(ctx context.Context, db qes.Querier, id int64) error {
	_, err := qes.ExecOne(ctx, db, dahua.Cameras.SELECT(dahua.Cameras.ID).WHERE(dahua.Cameras.ID.EQ(Int(id))))
	return err
}

// CameraGet [final]
func (dbT) CameraGet(ctx context.Context, db qes.Querier, id int64) (models.DahuaCamera, error) {
	var camera models.DahuaCamera
	err := qes.ScanOne(ctx, db, &camera, dahua.Cameras.
		SELECT(dbCameraProjection).
		WHERE(dahua.Cameras.ID.EQ(Int64(id))))
	return camera, err
}

// CameraGetByAddress [final]
func (dbT) CameraGetByAddress(ctx context.Context, db qes.Querier, address string) (models.DahuaCamera, error) {
	var camera models.DahuaCamera
	err := qes.ScanOne(ctx, db, &camera, dahua.Cameras.
		SELECT(dbCameraProjection).
		WHERE(dahua.Cameras.Address.EQ(String(address))))
	return camera, err
}

// CameraDelete [final]
func (dbT) CameraDelete(ctx context.Context, db qes.Querier, id int64) error {
	_, err := qes.ExecOne(ctx, db, dahua.Cameras.
		DELETE().
		WHERE(dahua.Cameras.ID.EQ(Int64(id))).
		RETURNING(dahua.Cameras.ID))
	return err
}

// CameraDetailUpdate [final]
func (dbT) CameraDetailUpdate(ctx context.Context, db qes.Querier, id int64, r models.DahuaCameraDetail) error {
	_, err := qes.ExecOne(ctx, db, dahua.CameraDetails.
		UPDATE(
			dahua.CameraDetails.Sn,
			dahua.CameraDetails.DeviceClass,
			dahua.CameraDetails.DeviceType,
			dahua.CameraDetails.HardwareVersion,
			dahua.CameraDetails.MarketArea,
			dahua.CameraDetails.ProcessInfo,
			dahua.CameraDetails.Vendor,
		).
		MODEL(model.CameraDetails{
			Sn:              r.SN,
			DeviceClass:     r.DeviceClass,
			DeviceType:      r.DeviceType,
			HardwareVersion: r.HardwareVersion,
			MarketArea:      r.MarketArea,
			ProcessInfo:     r.ProcessInfo,
			Vendor:          r.Vendor,
		}).
		WHERE(dahua.CameraDetails.CameraID.EQ(Int(id))))
	return err
}

// CameraSoftwaresUpdate [final]
func (dbT) CameraSoftwaresUpdate(ctx context.Context, db qes.Querier, id int64, r magicbox.GetSoftwareVersionResult) error {
	_, err := qes.ExecOne(ctx, db, dahua.CameraSoftwares.
		UPDATE(
			dahua.CameraSoftwares.Build,
			dahua.CameraSoftwares.BuildDate,
			dahua.CameraSoftwares.SecurityBaseLineVersion,
			dahua.CameraSoftwares.Version,
			dahua.CameraSoftwares.WebVersion,
		).
		MODEL(model.CameraSoftwares{
			Build:                   r.Build,
			BuildDate:               r.BuildDate,
			SecurityBaseLineVersion: r.SecurityBaseLineVersion,
			Version:                 r.Version,
			WebVersion:              r.WebVersion,
		}).
		WHERE(dahua.CameraSoftwares.CameraID.EQ(Int(id))))
	return err
}

// CameraLicensesReplace [final]
func (dbT) CameraLicensesReplace(ctx context.Context, db qes.Querier, id int64, licenses []license.LicenseInfo) error {
	return pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		_, err := qes.Exec(ctx, tx, dahua.CameraLicenses.DELETE().WHERE(dahua.CameraLicenses.CameraID.EQ(Int(id))))
		if err != nil {
			return err
		}

		stmt := dahua.CameraLicenses.
			INSERT(
				dahua.CameraLicenses.CameraID,
				dahua.CameraLicenses.AbroadInfo,
				dahua.CameraLicenses.AllType,
				dahua.CameraLicenses.DigitChannel,
				dahua.CameraLicenses.EffectiveDays,
				dahua.CameraLicenses.EffectiveTime,
				dahua.CameraLicenses.LicenseID,
				dahua.CameraLicenses.ProductType,
				dahua.CameraLicenses.Status,
				dahua.CameraLicenses.Username,
			)

		for _, r := range licenses {
			stmt = stmt.
				MODEL(model.CameraLicenses{
					CameraID:      int32(id),
					AbroadInfo:    r.AbroadInfo,
					AllType:       r.AllType,
					DigitChannel:  int32(r.DigitChannel),
					EffectiveDays: int32(r.EffectiveDays),
					EffectiveTime: int32(r.EffectiveTime),
					LicenseID:     int32(r.LicenseID),
					ProductType:   r.ProductType,
					Status:        int32(r.Status),
					Username:      r.Username,
				})
		}

		_, err = qes.ExecOne(ctx, tx, stmt)
		return err
	})
}

func (dbT) ScanCameraFilesUpsert(ctx context.Context, db qes.Querier, scannedAt time.Time, cam models.DahuaScanCursor, files []mediafilefind.FindNextFileInfo) (int64, error) {
	if len(files) == 0 {
		return 0, nil
	}

	// TODO: use MERGE instead of INSERT to prevent exhausting IDs
	stmt := dahua.ScanCameraFiles.
		INSERT(
			dahua.ScanCameraFiles.CameraID,
			dahua.ScanCameraFiles.FilePath,
			dahua.ScanCameraFiles.Kind,
			dahua.ScanCameraFiles.Size,
			dahua.ScanCameraFiles.StartTime,
			dahua.ScanCameraFiles.EndTime,
			dahua.ScanCameraFiles.Events,
			dahua.ScanCameraFiles.ScannedAt,
		)

	for _, file := range files {
		startTime, endTime, err := file.UniqueTime(cam.Seed, cam.Location.Location)
		if err != nil {
			return 0, err
		}

		// TODO: look into pgx's JSON type
		events, err := json.Marshal(file.Events)
		if err != nil {
			return 0, err
		}

		stmt = stmt.MODEL(
			model.ScanCameraFiles{
				CameraID:  int32(cam.CameraID),
				FilePath:  file.FilePath,
				Kind:      file.Type,
				Size:      int32(file.Length),
				StartTime: startTime,
				EndTime:   endTime,
				ScannedAt: scannedAt,
				Events:    string(events),
			},
		)
	}

	stmt = stmt.
		// FIXME: dahua.camera_files.start_time unique conflict should be handled somehow
		ON_CONFLICT(dahua.ScanCameraFiles.CameraID, dahua.ScanCameraFiles.FilePath).
		DO_UPDATE(SET(
			dahua.ScanCameraFiles.CameraID.SET(dahua.ScanCameraFiles.EXCLUDED.CameraID),
			dahua.ScanCameraFiles.FilePath.SET(dahua.ScanCameraFiles.EXCLUDED.FilePath),
			dahua.ScanCameraFiles.Kind.SET(dahua.ScanCameraFiles.EXCLUDED.Kind),
			dahua.ScanCameraFiles.Size.SET(dahua.ScanCameraFiles.EXCLUDED.Size),
			dahua.ScanCameraFiles.StartTime.SET(dahua.ScanCameraFiles.EXCLUDED.StartTime),
			dahua.ScanCameraFiles.EndTime.SET(dahua.ScanCameraFiles.EXCLUDED.EndTime),
			dahua.ScanCameraFiles.ScannedAt.SET(dahua.ScanCameraFiles.EXCLUDED.ScannedAt),
			dahua.ScanCameraFiles.Events.SET(dahua.ScanCameraFiles.EXCLUDED.Events)))

	res, err := qes.Exec(ctx, db, stmt)
	return res.RowsAffected(), err
}

// ScanCameraFilesDelete [final]
func (dbT) ScanCameraFilesDelete(ctx context.Context, db qes.Querier, scannedAt time.Time, cameraID int64, scanPeriod ScanPeriod) (int64, error) {
	res, err := qes.Exec(ctx, db, dahua.ScanCameraFiles.
		DELETE().
		WHERE(dahua.ScanCameraFiles.ScannedAt.LT(TimestampzT(scannedAt)).
			AND(dahua.ScanCameraFiles.CameraID.EQ(Int64(cameraID))).
			AND(dahua.ScanCameraFiles.StartTime.GT_EQ(TimestampzT(scanPeriod.Start))).
			AND(dahua.ScanCameraFiles.StartTime.LT(TimestampzT(scanPeriod.End)))))
	return res.RowsAffected(), err
}

func (dbT) ScanCursorGet(ctx context.Context, db qes.Querier, cameraID int64) (models.DahuaScanCursor, error) {
	var res models.DahuaScanCursor
	err := qes.ScanOne(ctx, db, &res,
		SELECT(
			dahua.Cameras.ID.AS("camera_id"),
			Raw(fmt.Sprintf("coalesce(%s, %s.%s)", dahua.ScanSeeds.Seed.Name(), dahua.Cameras.ID.TableName(), dahua.Cameras.ID.Name())).AS(dahua.ScanSeeds.Seed.Name()),
			dahua.Cameras.Location.AS("location"),
			dahua.ScanCursors.FullComplete.AS("full_complete"),
			dahua.ScanCursors.FullCursor.AS("full_cursor"),
			dahua.ScanCursors.FullEpoch.AS("full_epoch"),
			dahua.ScanCursors.FullEpochEnd.AS("full_epoch_end"),
			dahua.ScanCursors.QuickCursor.AS("quick_cursor"),
		).FROM(dahua.Cameras.
			LEFT_JOIN(dahua.ScanCursors, dahua.ScanCursors.CameraID.EQ(dahua.Cameras.ID)).
			LEFT_JOIN(dahua.ScanSeeds, dahua.ScanSeeds.CameraID.EQ(dahua.Cameras.ID))).
			WHERE(dahua.Cameras.ID.EQ(Int(cameraID))))
	return res, err
}

func (dbT) ScanCursorReset(ctx context.Context, db qes.Querier, cameraID int64) error {
	_, err := qes.ExecOne(ctx, db, dahua.ScanCursors.
		UPDATE().
		SET(
			dahua.ScanCursors.QuickCursor.SET(RawTimestampz("default")),
			dahua.ScanCursors.FullCursor.SET(RawTimestampz("default")),
			dahua.ScanCursors.FullEpochEnd.SET(RawTimestampz("default")),
		).
		WHERE(dahua.ScanCursors.CameraID.EQ(Int(cameraID))))
	return err
}

// ---------- ScanCursorLock

type ScanCursorLock struct {
	models.DahuaScanCursor
	tx qes.Querier
}

func (dbT) newScanCursorLock(tx qes.Querier, scanCursor models.DahuaScanCursor) ScanCursorLock {
	return ScanCursorLock{
		DahuaScanCursor: scanCursor,
		tx:              tx,
	}
}

func (s ScanCursorLock) UpdateFullCursor(ctx context.Context, fullCursor time.Time) error {
	_, err := qes.ExecOne(ctx, s.tx, dahua.ScanCursors.
		UPDATE().
		SET(dahua.ScanCursors.FullCursor.SET(TimestampzT(fullCursor))).
		WHERE(dahua.ScanCursors.CameraID.EQ(Int(s.DahuaScanCursor.CameraID))))
	return err
}

func (s ScanCursorLock) UpdateFullCursorFromActiveScanTaskCursor(ctx context.Context) error {
	_, err := qes.ExecOne(ctx, s.tx, dahua.ScanCursors.
		UPDATE(dahua.ScanCursors.FullCursor).
		SET(
			dahua.ScanActiveTasks.SELECT(dahua.ScanActiveTasks.Cursor).WHERE(dahua.ScanActiveTasks.CameraID.EQ(Int(s.DahuaScanCursor.CameraID))),
		).
		WHERE(dahua.ScanCursors.CameraID.EQ(Int(s.DahuaScanCursor.CameraID))))
	return err
}

func (s ScanCursorLock) UpdateQuickCursor(ctx context.Context, quickCursor time.Time) error {
	_, err := qes.ExecOne(ctx, s.tx, dahua.ScanCursors.
		UPDATE().
		SET(dahua.ScanCursors.QuickCursor.SET(TimestampzT(quickCursor))).
		WHERE(dahua.ScanCursors.CameraID.EQ(Int(s.DahuaScanCursor.CameraID))),
	)
	return err
}

// ----------

// ScanQueueTaskCreate [final]
func (dbT) ScanQueueTaskCreate(ctx context.Context, db qes.Querier, r models.DahuaScanQueueTask) error {
	return pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		res, err := qes.Exec(ctx, tx, dahua.ScanQueueTasks.
			SELECT(dahua.ScanQueueTasks.ID).
			WHERE(AND(
				dahua.ScanQueueTasks.CameraID.EQ(Int(r.CameraID)),
				dahua.ScanQueueTasks.Kind.EQ(NewEnumValue(r.Kind.String()))),
			).
			FOR(UPDATE()))
		if err != nil {
			return err
		}
		if err := ScanKindQuota(r.Kind, res.RowsAffected()); err != nil {
			return err
		}

		_, err = qes.Exec(ctx, tx, dahua.ScanQueueTasks.
			INSERT(
				dahua.ScanQueueTasks.CameraID,
				dahua.ScanQueueTasks.Kind,
				dahua.ScanQueueTasks.Range,
			).
			MODEL(struct {
				model.ScanQueueTasks
				Range models.DahuaScanRange
			}{
				ScanQueueTasks: model.ScanQueueTasks{
					CameraID: int32(r.CameraID),
					Kind:     r.Kind,
				},
				Range: r.Range,
			}))
		return err
	})
}

// ScanQueueTaskNext pops a queued task off the queue and locks the scan cursor.
// WARNING: fn can only access scan tables or else it will deadlock when the camera gets deleted.
func (dbT) ScanQueueTaskNext(ctx context.Context, db qes.Querier, fn func(ctx context.Context, scanCursorLock ScanCursorLock, queueTask models.DahuaScanQueueTask) error) (bool, error) {
	var ok bool
	return ok, pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		// Lock scan cursor and get the first queued scan task
		var queueTask models.DahuaScanQueueTask
		err := qes.ScanOne(ctx, tx, &queueTask, RawStatement(`
			SELECT
				qt.id,
				qt.camera_id,
				qt.kind,
				qt."range"
			FROM dahua.scan_cursors sc
			LEFT JOIN dahua.scan_queue_tasks qt ON qt.camera_id = sc.camera_id
			WHERE (qt.id notnull)
			ORDER BY qt.id
			LIMIT 1
			FOR NO KEY UPDATE OF sc SKIP LOCKED
		`)) // NOTE: Using raw sql because jet does not support OF clause for row locks
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return err
		}
		ok = true

		// Get scan cursor lock
		scanCursor, err := DB.ScanCursorGet(ctx, tx, queueTask.CameraID)
		if err != nil {
			return err
		}
		scanCursorLock := DB.newScanCursorLock(tx, scanCursor)

		if err := fn(ctx, scanCursorLock, queueTask); err != nil {
			return err
		}

		// Delete the queued scan task
		_, err = qes.Exec(ctx, tx, dahua.ScanQueueTasks.DELETE().WHERE(dahua.ScanQueueTasks.ID.EQ(Int(queueTask.ID))))
		return err
	})
}

// ScanActiveQueueClear clears orphan active scans.
func (dbT) ScanActiveQueueClear(ctx context.Context, db qes.Querier) error {
	_, err := qes.Exec(ctx, db, dahua.ScanActiveTasks.
		DELETE().
		WHERE(dahua.ScanActiveTasks.CameraID.IN(
			dahua.ScanCursors.SELECT(dahua.ScanCursors.CameraID).FOR(UPDATE().SKIP_LOCKED()))))
	return err
}

func (dbT) ScanActiveTaskCreate(ctx context.Context, db qes.Querier, r models.DahuaScanQueueTask) (models.DahuaScanActiveTask, error) {
	var res models.DahuaScanActiveTask
	// TODO: instead of deleting the active scan if it exists, it should return the active scan instead so we can have continue scanning on restart
	err := pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		_, err := qes.Exec(ctx, tx, dahua.ScanActiveTasks.
			DELETE().
			WHERE(dahua.ScanActiveTasks.QueueID.EQ(Int(r.ID))))
		if err != nil {
			return err
		}

		err = qes.ScanOne(ctx, tx, &res, dahua.ScanActiveTasks.
			INSERT(
				dahua.ScanActiveTasks.CameraID,
				dahua.ScanActiveTasks.QueueID,
				dahua.ScanActiveTasks.Kind,
				dahua.ScanActiveTasks.Range,
				dahua.ScanActiveTasks.Cursor,
			).
			MODEL(struct {
				model.ScanActiveTasks
				Range models.DahuaScanRange
			}{
				ScanActiveTasks: model.ScanActiveTasks{
					CameraID: int32(r.CameraID),
					QueueID:  int32(r.ID),
					Kind:     r.Kind,
					Cursor:   r.Range.End,
				},
				Range: r.Range,
			}).
			RETURNING(
				dahua.ScanActiveTasks.CameraID.AS("camera_id"),
				dahua.ScanActiveTasks.Kind.AS("kind"),
				dahua.ScanActiveTasks.Range.AS("range"),
				dahua.ScanActiveTasks.Cursor.AS("cursor"),
				dahua.ScanActiveTasks.StartedAt.AS("started_at"),
				dahua.ScanActiveTasks.Deleted.AS("deleted"),
				dahua.ScanActiveTasks.Upserted.AS("upserted"),
				dahua.ScanActiveTasks.Percent.AS("percent")))

		return err
	})
	return res, err
}

func (dbT) ScanActiveTaskComplete(ctx context.Context, db qes.Querier, cameraID int64, errString string) error {
	return pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		_, err := qes.ExecOne(ctx, tx, dahua.ScanCompleteTasks.
			INSERT(
				dahua.ScanCompleteTasks.CameraID,
				dahua.ScanCompleteTasks.Kind,
				dahua.ScanCompleteTasks.Range,
				dahua.ScanCompleteTasks.Cursor,
				dahua.ScanCompleteTasks.StartedAt,
				dahua.ScanCompleteTasks.Deleted,
				dahua.ScanCompleteTasks.Upserted,
				dahua.ScanCompleteTasks.Percent,
				dahua.ScanCompleteTasks.Duration,
				dahua.ScanCompleteTasks.Error,
			).
			QUERY(dahua.ScanActiveTasks.
				SELECT(
					dahua.ScanActiveTasks.CameraID,
					dahua.ScanActiveTasks.Kind,
					dahua.ScanActiveTasks.Range,
					dahua.ScanActiveTasks.Cursor,
					dahua.ScanActiveTasks.StartedAt,
					dahua.ScanActiveTasks.Deleted,
					dahua.ScanActiveTasks.Upserted,
					dahua.ScanActiveTasks.Percent,
					Raw(fmt.Sprintf("EXTRACT(EPOCH FROM age(CURRENT_TIMESTAMP, %s.%s))", dahua.ScanActiveTasks.StartedAt.TableName(), dahua.ScanActiveTasks.StartedAt.Name())),
					String(errString).AS("error"),
				).
				WHERE(dahua.ScanActiveTasks.CameraID.EQ(Int(cameraID)))))
		if err != nil {
			return err
		}

		_, err = qes.ExecOne(ctx, tx, dahua.ScanActiveTasks.
			DELETE().
			WHERE(dahua.ScanActiveTasks.CameraID.EQ(Int(cameraID))))

		return err
	})
}

// ScanActiveProgressUpdate [final]
func (dbT) ScanActiveProgressUpdate(ctx context.Context, db qes.Querier, r models.DahuaScanActiveProgress) (models.DahuaScanActiveProgress, error) {
	var res models.DahuaScanActiveProgress
	err := qes.ScanOne(ctx, db, &res, dahua.ScanActiveTasks.
		UPDATE(
			dahua.ScanActiveTasks.Upserted,
			dahua.ScanActiveTasks.Deleted,
			dahua.ScanActiveTasks.Percent,
			dahua.ScanActiveTasks.Cursor,
		).
		MODEL(model.ScanActiveTasks{
			Upserted: int32(r.Upserted),
			Deleted:  int32(r.Deleted),
			Percent:  float32(r.Percent),
			Cursor:   r.Cursor,
		}).
		WHERE(dahua.ScanActiveTasks.CameraID.EQ(Int(r.CameraID))).
		RETURNING(
			dahua.ScanActiveTasks.CameraID.AS("camera_id"),
			dahua.ScanActiveTasks.Upserted.AS("upserted"),
			dahua.ScanActiveTasks.Deleted.AS("deleted"),
			dahua.ScanActiveTasks.Percent.AS("percent"),
			dahua.ScanActiveTasks.Cursor.AS("cursor")))
	return res, err
}
