package wrap

import (
	"context"

	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/pg"

	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/wrap"
)

func (r WrapRepository) AddWrap(ctx context.Context, wrap wrap.Wrap) error {
	tx, err := pg.GetTransactionFromContext(ctx)
	if err != nil {
		return err
	}

	const query = `INSERT INTO wrap(name, max_weight, cost) VALUES ($1, $2, $3)`

	_, err = tx.Exec(
		ctx,
		query,
		wrap.Name,
		wrap.MaxWeight,
		wrap.Cost,
	)

	return err
}
