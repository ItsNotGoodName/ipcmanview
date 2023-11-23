package sqlc

import (
	"context"
	"database/sql"
	"errors"
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
	err = tx.setDahuaSeed(ctx, sql.NullInt64{
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

// func (db DB) UpsertDahuaCamera(ctx context.Context, id int64, args CreateDahuaCameraParams) (int64, error) {
// 	_, err := db.UpdateDahuaCamera(ctx, UpdateDahuaCameraParams{
// 		Name:     args.Name,
// 		Address:  args.Address,
// 		Username: args.Username,
// 		Password: args.Password,
// 		Location: args.Location,
// 		ID:       id,
// 	})
// 	if err == nil {
// 		return id, nil
// 	}
// 	if !errors.Is(err, sql.ErrNoRows) {
// 		return 0, err
// 	}
//
// 	return db.createDahuaCamera(ctx, args)
// }

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
