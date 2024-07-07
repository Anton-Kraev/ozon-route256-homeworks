package cli

import (
	"context"
	"strings"

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

			return c.txMiddleware.CreateTransactionContext(ctx, txWriteOptions, args, c.wrapController.AddWrap)
		})
	default:
		cancel()
	}
}
