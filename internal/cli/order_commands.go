package cli

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"

	domain "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/event"
)

func (c *CLI) handleOrderCommand(input []string) {
	var (
		inputString = strings.Join(input, " ")
		cmd         = commandName(input[0])
		args        = input[1:]

		txReadOptions  = pgx.TxOptions{AccessMode: pgx.ReadOnly, IsoLevel: pgx.ReadCommitted}
		txWriteOptions = pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead}

		ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	)

	c.cmdCounter++

	addTask := func(txOptions pgx.TxOptions, handler func(context.Context, []string) (string, error)) {
		c.workerPool.AddTask(c.cmdCounter, inputString, func() (string, error) {
			defer cancel()

			return c.txMiddleware.CreateTransactionContext(ctx, txOptions, args,
				func(ctx context.Context, args []string) (string, error) {
					err := c.eventRepo.AddEvent(ctx, domain.Event{
						Type:    domain.Type(cmd),
						Payload: inputString,
					})
					if err != nil {
						return "", err
					}

					return handler(ctx, args)
				},
			)
		})
	}

	switch cmd {
	case receiveOrder:
		addTask(txWriteOptions, c.orderController.ReceiveOrder)
	case returnOrder:
		addTask(txWriteOptions, c.orderController.ReturnOrder)
	case deliverOrders:
		addTask(txWriteOptions, c.orderController.DeliverOrders)
	case clientOrders:
		addTask(txReadOptions, c.orderController.ClientOrders)
	case refundOrder:
		addTask(txWriteOptions, c.orderController.RefundOrder)
	case refundsList:
		addTask(txReadOptions, c.orderController.RefundsList)
	default:
		cancel()
	}
}
