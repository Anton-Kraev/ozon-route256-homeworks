package wrap

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"

	"gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/wrap"
)

func (r WrapRepository) GetWrapByName(ctx context.Context, name string) (*wrap.Wrap, error) {
	const query = `SELECT * FROM wrap WHERE name = $1`

	var wraps []WrapSchema

	err := pgxscan.Select(ctx, r.pool, &wraps, query, name)
	if err != nil || len(wraps) == 0 {
		return nil, err
	}

	wrapDomain := wraps[0].ToDomain()

	return &wrapDomain, nil
}
