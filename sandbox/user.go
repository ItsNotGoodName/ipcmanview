package sandbox

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/ipcmango/internal/core"
	"github.com/ItsNotGoodName/ipcmango/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func User(ctx context.Context, pool *pgxpool.Pool) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get connection")
	}
	defer conn.Release()

	user, err := core.NewUser(core.UserCreate{
		Email:           "admin@example.com",
		Username:        "admin",
		Password:        "password",
		PasswordConfirm: "password",
	})
	if err != nil {
		log.Err(err).Msg("Failed to create user")
		return
	}

	_, err = db.UserCreate(db.Context{
		Context: ctx,
		Conn:    conn.Conn(),
	}, user)
	if err != nil {
		log.Err(err).Msg("Failed to persist user")
	}

	err = core.UserCheckPassword(user, "password")
	fmt.Println(err)
}
