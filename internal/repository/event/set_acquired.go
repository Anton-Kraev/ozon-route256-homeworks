package event

import (
	"context"
	"time"

	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/pg"
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
