package event

import (
	"context"
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-6/internal/pg"
)

func (r EventRepository) SetAcquired(ctx context.Context, eventID uint64, acquiredTo time.Time) error {
	tx, err := pg.GetTransactionFromContext(ctx)
	if err != nil {
		return err
	}

	const query = `UPDATE event SET acquired_to = $2 WHERE id = $1`

	_, err = tx.Exec(ctx, query, eventID, acquiredTo)

	return err
}
