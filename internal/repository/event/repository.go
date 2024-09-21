package event

import "github.com/jackc/pgx/v5/pgxpool"

type EventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) EventRepository {
	return EventRepository{pool: pool}
}
