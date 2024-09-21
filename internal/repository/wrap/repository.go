package wrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type WrapRepository struct {
	pool *pgxpool.Pool
}

func NewWrapRepository(pool *pgxpool.Pool) WrapRepository {
	return WrapRepository{pool: pool}
}
