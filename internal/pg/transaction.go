package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type txKey struct{}

func GetTransactionFromContext(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if !ok {
		return nil, errors.New("can't pick transaction from context")
	}

	return tx, nil
}

func AddTransactionToContext(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}
