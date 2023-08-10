package dbtest

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/ItsNotGoodName/ipcmanview/migrations"
	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context) (*pgx.Conn, func()) {
	url := "postgres://postgres:postgres@localhost:5432"
	database := "postgres_test_" + strconv.Itoa(rand.Int())

	// ---------------------- Initialize database

	tempConn, err := pgx.Connect(ctx, url)
	if err != nil {
		panic(err)
	}

	_, err = tempConn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %s`, database))
	tempConn.Close(ctx)
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

	return conn, close
}
