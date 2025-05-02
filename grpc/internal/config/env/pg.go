package env

import (
	"errors"
	"github.com/beachrockhotel/internal/config"
	"os"
)

var _ config.PgConfig = (*pgConfig)(nil)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		errors.New("config env var " + dsnEnvName + " is not set")
	}
	return &pgConfig{dsn: dsn}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
