package config

import (
	"flag"
)

type Web struct {
	HTTPHost string
	HTTPPort string
	DBPath   string
}

func NewWeb(flags *flag.FlagSet) *Web {
	c := &Web{
		HTTPHost: getEnvOptional("HTTP_HOST", ""),
		HTTPPort: getEnvOptional("HTTP_PORT", "8080"),
		DBPath:   getEnvOptional("DB_PATH", "sqlite.db"),
	}
	flags.StringVar(&c.HTTPHost, "http-host", c.HTTPHost, "HTTP host to listen on.")
	flags.StringVar(&c.HTTPPort, "http-port", c.HTTPPort, "HTTP port to listen on.")
	flags.StringVar(&c.DBPath, "db-path", c.DBPath, "Database file path.")
	return c
}
