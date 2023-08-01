package db

import (
	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/internal/db/gen/postgres/public/model"
	. "github.com/ItsNotGoodName/ipcmango/internal/db/gen/postgres/public/table"
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

func UserCreate(ctx Context, r core.User) (core.User, error) {
	var user core.User
	err := qes.ScanOne(ctx.Context, ctx.Conn, &user, Users.
		INSERT(Users.Email, Users.Username, Users.Password).
		MODEL(model.Users{
			Email:    r.Email,
			Username: r.Username,
			Password: r.Password,
		}).
		RETURNING(pUsers),
	)
	return user, err
}

func UserGetByUsernameOrEmail(ctx Context, usernameOrEmail string) (core.User, error) {
	var user core.User
	err := qes.ScanOne(ctx.Context, ctx.Conn, &user, Users.
		SELECT(pUsers).
		WHERE(
			Users.Email.EQ(String(usernameOrEmail)).OR(Users.Username.EQ(String(usernameOrEmail))),
		))
	return user, err
}

func UserGet(ctx Context, id int64) (core.User, error) {
	var user core.User
	err := qes.ScanOne(ctx.Context, ctx.Conn, &user, Users.
		SELECT(pUsers).
		WHERE(
			Users.ID.EQ(Int64(id)),
		))
	return user, err
}
