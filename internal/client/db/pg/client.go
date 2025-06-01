package pg

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pkg/errors"

	"github.com/beachrockhotel/auth/internal/client/db"
)

type pgClient struct {
	masterDBC db.DB
}

func New(ctx context.Context, dsn string) (db.Client, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse config")
	}

	dbc, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pool")
	}

	return &pgClient{
		masterDBC: &pg{dbc: dbc},
	}, nil
}

func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
