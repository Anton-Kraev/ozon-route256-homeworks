package middlewares

import (
	"context"

	"github.com/jackc/pgx/v5"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/pg"
)

type PoolTxBeginner interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type TransactionMiddleware struct {
	pool PoolTxBeginner
}

func NewTransactionMiddleware(pool PoolTxBeginner) *TransactionMiddleware {
	return &TransactionMiddleware{pool: pool}
}

func (tm *TransactionMiddleware) CreateTransactionContext(
	ctx context.Context,
	txOptions pgx.TxOptions,
	args []string,
	handler func(ctx context.Context, args []string) (string, error),
) (res string, err error) {
	var tx pgx.Tx

	tx, err = tm.pool.BeginTx(ctx, txOptions)
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)

			panic(p)
		}

		if err != nil {
			_ = tx.Rollback(ctx)

			return
		}

		err = tx.Commit(ctx)
	}()

	ctx = pg.AddTransactionToContext(ctx, tx)
	res, err = handler(ctx, args)

	return
}
