package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPoolConn(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	return pgxpool.NewWithConfig(ctx, pgxConfig)
}
