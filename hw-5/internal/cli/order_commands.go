package cli

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
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

			return c.txMiddleware.CreateTransactionContext(ctx, txOptions, args, handler)
		})
	}

	switch cmd {
	case receiveOrder:
		addTask(txWriteOptions, c.controller.ReceiveOrder)
	case returnOrder:
		addTask(txWriteOptions, c.controller.ReturnOrder)
	case deliverOrders:
		addTask(txWriteOptions, c.controller.DeliverOrders)
	case clientOrders:
		addTask(txReadOptions, c.controller.ClientOrders)
	case refundOrder:
		addTask(txWriteOptions, c.controller.RefundOrder)
	case refundsList:
		addTask(txReadOptions, c.controller.RefundsList)
	default:
		cancel()
	}
}
