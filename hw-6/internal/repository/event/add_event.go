package event

import (
	"context"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-6/internal/models/domain/event"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-6/internal/pg"
)

func (r EventRepository) AddEvent(ctx context.Context, event event.Event) error {
	tx, err := pg.GetTransactionFromContext(ctx)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO event(id, type, payload, processed_at, acquired_to) 
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = tx.Exec(
		ctx,
		query,
		event.ID,
		event.Type,
		event.Payload,
		event.ProcessedAt,
		event.AcquiredTo,
	)

	return err
}
