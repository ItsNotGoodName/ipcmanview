package sqlc

import (
	"context"
	"database/sql"
	"errors"
)

type CreateDahuaCameraParams = createDahuaCameraParams

func (db DB) CreateDahuaCamera(ctx context.Context, arg CreateDahuaCameraParams) (int64, error) {
	tx, err := db.BeginTx(ctx, true)
	if err != nil {
		return 0, nil
	}
	defer tx.Rollback()

	id, err := db.createDahuaCamera(ctx, arg)
	if err != nil {
		return id, err
	}

	// TODO: sql.NullInt64 should just be int64...
	err = db.setDahuaSeed(ctx, sql.NullInt64{
		Valid: true,
		Int64: id,
	})
	if err != nil {
		return id, err
	}

	err = tx.Commit()
	if err != nil {
		return id, err
	}

	return id, nil
}

func (db DB) UpsertDahuaCamera(ctx context.Context, id int64, args CreateDahuaCameraParams) (int64, error) {
	_, err := db.UpdateDahuaCamera(ctx, UpdateDahuaCameraParams{
		Name:     args.Name,
		Address:  args.Address,
		Username: args.Username,
		Password: args.Password,
		Location: args.Location,
		ID:       id,
	})
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	return db.createDahuaCamera(ctx, args)
}
