package db

import (
	"context"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/internal/dbgen/postgres/public/model"
	. "github.com/ItsNotGoodName/ipcmango/internal/dbgen/postgres/public/table"
	"github.com/ItsNotGoodName/ipcmango/pkg/qes"
	. "github.com/go-jet/jet/v2/postgres"
)

var pUsers ProjectionList = []Projection{
	Users.ID.AS("id"),
	Users.Email.AS("email"),
	Users.Username.AS("username"),
	Users.Password.AS("password"),
	Users.CreatedAt.AS("created_at"),
}

type userT struct{}

var User userT

func (userT) Create(ctx context.Context, db qes.Querier, u core.User) (core.User, error) {
	var user core.User
	err := qes.ScanOne(ctx, db, &user, Users.
		INSERT(Users.Email, Users.Username, Users.Password).
		MODEL(model.Users{
			Email:    u.Email,
			Username: u.Username,
			Password: u.Password,
		}).
		RETURNING(pUsers),
	)
	return user, err
}

func (userT) GetByUsernameOrEmail(ctx context.Context, db qes.Querier, usernameOrEmail string) (core.User, error) {
	var user core.User
	err := qes.ScanOne(ctx, db, &user, Users.
		SELECT(pUsers).
		WHERE(
			Users.Email.EQ(String(usernameOrEmail)).OR(Users.Username.EQ(String(usernameOrEmail))),
		))
	return user, err
}

func (userT) Get(ctx context.Context, db qes.Querier, id int64) (core.User, error) {
	var user core.User
	err := qes.ScanOne(ctx, db, &user, Users.
		SELECT(pUsers).
		WHERE(
			Users.ID.EQ(Int64(id)),
		))
	return user, err
}
