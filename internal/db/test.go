package db

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/ipcmango/migrations"
	"github.com/jackc/pgx/v5"
)

// TestConnect is only used for testing.
func TestConnect(ctx context.Context) (Context, func()) {
	url := "postgres://postgres:postgres@localhost:5432"
	database := "postgres_test"

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, database))
	if err != nil {
		conn.Close(ctx)
		panic(err)
	}

	_, err = conn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %s`, database))
	conn.Close(ctx)
	if err != nil {
		panic(err)
	}

	conn, err = pgx.Connect(ctx, url+"/"+database)
	if err != nil {
		panic(err)
	}

	err = migrations.MigrateConn(ctx, conn)
	if err != nil {
		conn.Close(ctx)
		panic(err)
	}

	return Context{
			Context: ctx,
			Conn:    conn,
		}, func() {
			conn.Close(ctx)
		}
}
