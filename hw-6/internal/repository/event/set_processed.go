package event

import (
	"context"
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-6/internal/pg"
)

func (r EventRepository) SetProcessed(ctx context.Context, eventID uint64, processedAt time.Time) error {
	tx, err := pg.GetTransactionFromContext(ctx)
	if err != nil {
		return err
	}

	const query = `UPDATE event SET processed_at = $2 WHERE id = $1`

	_, err = tx.Exec(ctx, query, eventID, processedAt)

	return err
}
