package outbox

import (
	"context"
	"fmt"
	"time"

	domain "gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/event"
)

type eventRepo interface {
	AddEvent(ctx context.Context, event domain.Event) error
	GetNewEvent(ctx context.Context) (*domain.Event, error)
	SetProcessed(ctx context.Context, eventID uint64, processedAt time.Time) error
	SetAcquired(ctx context.Context, eventID uint64, acquiredTo time.Time) error
}

type OutboxProcessor struct {
	eventRepo eventRepo
}

func NewOutboxProcessor(eventRepo eventRepo) *OutboxProcessor {
	return &OutboxProcessor{eventRepo: eventRepo}
}

func (p *OutboxProcessor) Start(ctx context.Context, handlePeriod time.Duration) {
	ticker := time.NewTicker(handlePeriod)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping processing events")

			return
		case <-ticker.C:
		}

		event, err := p.eventRepo.GetNewEvent(ctx)
		if err != nil {
			fmt.Println("Error getting new event:", err)

			continue
		}

		if event == nil {
			continue
		}

		// TODO: send event

		if err = p.eventRepo.SetProcessed(ctx, event.ID, time.Now()); err != nil {
			fmt.Println(err)
		}
	}
}
