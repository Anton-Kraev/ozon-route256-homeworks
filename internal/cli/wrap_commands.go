package cli

import (
	"context"
	"strings"

	domain "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/event"

	"github.com/jackc/pgx/v5"
)

func (c *CLI) handleWrapCommand(input []string) {
	var (
		inputString = strings.Join(input, " ")
		cmd         = commandName(input[0])
		args        = input[1:]

		txWriteOptions = pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead}

		ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	)

	switch cmd {
	case addWrap:
		c.cmdCounter++
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			defer cancel()

			return c.txMiddleware.CreateTransactionContext(ctx, txWriteOptions, args,
				func(ctx context.Context, args []string) (string, error) {
					err := c.eventRepo.AddEvent(ctx, domain.Event{
						Type:    domain.Type(cmd),
						Payload: inputString,
					})
					if err != nil {
						return "", err
					}

					return c.wrapController.AddWrap(ctx, args)
				},
			)
		})
	default:
		cancel()
	}
}
