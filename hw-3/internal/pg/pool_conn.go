package pg

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPoolConn(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := os.Getenv("PGDATABASE")

	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	return pgxpool.NewWithConfig(ctx, pgxConfig)
}
