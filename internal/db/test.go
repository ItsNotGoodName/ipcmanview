package db

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/ItsNotGoodName/ipcmango/migrations"
	"github.com/jackc/pgx/v5"
)

// TestConnect is only used for testing.
func TestConnect(ctx context.Context) (Context, func()) {
	url := "postgres://postgres:postgres@localhost:5432"
	database := "postgres_test_" + strconv.Itoa(rand.Int())

	// ---------------------- Initialize database

	initConn, err := pgx.Connect(ctx, url)
	if err != nil {
		panic(err)
	}

	_, err = initConn.Exec(ctx, fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, database))
	if err != nil {
		initConn.Close(ctx)
		panic(err)
	}

	_, err = initConn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %s`, database))
	initConn.Close(ctx)
	if err != nil {
		panic(err)
	}

	// ---------------------- Connect

	conn, err := pgx.Connect(ctx, url+"/"+database)
	if err != nil {
		panic(err)
	}
	close := func() {
		conn.Close(ctx)

		conn, err := pgx.Connect(ctx, url)
		if err != nil {
			return
		}

		conn.Exec(ctx, fmt.Sprintf(`DROP DATABASE %s`, database))

		conn.Close(ctx)
	}

	err = migrations.MigrateConn(ctx, conn)
	if err != nil {
		conn.Close(ctx)
		panic(err)
	}

	return Context{
		Context: ctx,
		Conn:    conn,
	}, close
}
