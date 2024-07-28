package event

import (
	"time"
)

type Event struct {
	ID          uint64
	Type        Type
	Payload     string
	ProcessedAt time.Time
	AcquiredTo  time.Time
}
