package event

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"

	models "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/event"
)

func (r EventRepository) GetNewEvent(ctx context.Context) (*models.Event, error) {
	const query = `
		SELECT * FROM event 
		WHERE processed_at IS NULL AND (acquired_to IS NULL OR acquired_to < NOW()) 
		LIMIT 1
	`

	var events []EventSchema
	if err := pgxscan.Select(ctx, r.pool, &events, query); err != nil || len(events) == 0 {
		return nil, err
	}

	event := events[0].ToDomain()

	return &event, nil
}
