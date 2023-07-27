package sandbox

import (
	"context"
	"fmt"

	"github.com/ItsNotGoodName/ipcmango/internal/db"
	. "github.com/ItsNotGoodName/ipcmango/internal/db/gen/postgres/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func Jet(ctx context.Context, pool *pgxpool.Pool) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get connection")
	}
	defer conn.Release()

	_, err = db.Exec(ctx, conn, Placeholder.INSERT(Placeholder.Title).VALUES("Hello"))
	if err != nil {
		log.Err(err).Msg("Failed to insert into placeholder")
		return
	}

	type placeholder struct {
		ID    string
		Title string
	}

	pl := []placeholder{}
	err = db.Scan(ctx, conn, &pl, SELECT(Placeholder.ID.AS("id"), Placeholder.Title.AS("title")).FROM(Placeholder))
	if err != nil {
		log.Err(err).Msg("Failed to list placeholders")
		return
	}

	for _, p := range pl {
		fmt.Printf("p: %+v\n", p)

	}
}
