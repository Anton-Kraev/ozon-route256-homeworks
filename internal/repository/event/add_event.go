package event

import (
	"context"

	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/event"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/pg"
)

func (r EventRepository) AddEvent(ctx context.Context, event event.Event) error {
	tx, err := pg.GetTransactionFromContext(ctx)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO event(type, payload, processed_at, acquired_to) 
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(
		ctx,
		query,
		event.Type,
		event.Payload,
		event.ProcessedAt,
		event.AcquiredTo,
	)

	return err
}
