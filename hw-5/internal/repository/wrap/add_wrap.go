package wrap

import (
	"context"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/pg"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
)

func (r WrapRepository) AddWrap(ctx context.Context, wrap wrap.Wrap) error {
	tx, err := pg.GetTransactionFromContext(ctx)
	if err != nil {
		return err
	}

	const query = `INSERT INTO wrap(name, weight, cost) VALUES ($1, $2, $3)`

	_, err = tx.Exec(
		ctx,
		query,
		wrap.Name,
		wrap.Weight,
		wrap.Cost,
	)

	return err
}
