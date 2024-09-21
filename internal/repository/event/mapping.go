package event

import (
	"database/sql"
	"time"

	"gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/event"
)

type EventSchema struct {
	ID          uint64       `db:"id"`
	Type        string       `db:"type"`
	Payload     string       `db:"payload"`
	ProcessedAt sql.NullTime `db:"processed_at"`
	AcquiredTo  sql.NullTime `db:"acquired_to"`
	CreatedAt   time.Time    `db:"created_at"`
}

func (r EventSchema) ToDomain() event.Event {
	return event.Event{
		ID:          r.ID,
		Type:        event.Type(r.Type),
		Payload:     r.Payload,
		ProcessedAt: r.ProcessedAt.Time,
		AcquiredTo:  r.AcquiredTo.Time,
	}
}
