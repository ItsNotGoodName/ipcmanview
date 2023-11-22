package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"strings"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

//go:embed sql/*.sql
var migrations embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(db, "sql"); err != nil {
		return err
	}

	return nil
}

func init() {
	goose.SetLogger(&logger{})
}

type logger struct{}

func (*logger) Fatalf(format string, v ...interface{}) {
	log.Fatal().Msg(strings.TrimSuffix(fmt.Sprintf(format, v...), "\n"))
}

func (*logger) Printf(format string, v ...interface{}) {
	log.Info().Msg(strings.TrimSuffix(fmt.Sprintf(format, v...), "\n"))
}
