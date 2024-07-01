package wrap

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/middlewares"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
	wrapRepo "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/repository/wrap"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/tests/pg"
)

func TestAddWrap(t *testing.T) {
	var (
		ctx = context.Background()

		testWraps = []wrap.Wrap{
			{
				Name:   "box",
				Weight: 1000,
				Cost:   1000,
			},
			{
				Name:   "package",
				Weight: 2000,
				Cost:   2000,
			},
		}
	)

	pg.DB.SetUp(t, "orders", "wrap")
	defer pg.DB.TearDown(t)
	pg.DB.FillWraps(testWraps[:1])
	repo := wrapRepo.NewWrapRepository(pg.DB.ConnPool)
	txMw := middlewares.NewTransactionMiddleware(pg.DB.ConnPool)

	t.Run("add_existing", func(t *testing.T) {
		t.Parallel()

		_, err := txMw.CreateTransactionContext(
			ctx,
			pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead},
			[]string{},
			func(ctx context.Context, args []string) (string, error) {
				return "", repo.AddWrap(ctx, testWraps[0])
			},
		)

		require.Error(t, err)
	})

	t.Run("add_new", func(t *testing.T) {
		t.Parallel()

		res, err := txMw.CreateTransactionContext(
			ctx,
			pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead},
			[]string{},
			func(ctx context.Context, args []string) (string, error) {
				return "", repo.AddWrap(ctx, testWraps[1])
			},
		)

		require.NoError(t, err)
		require.Empty(t, res)

		wraps := pg.DB.GetAllWraps()
		require.Len(t, wraps, 2)
		AssertEqualWraps(t, testWraps[0], wraps[0].ToDomain())
		AssertEqualWraps(t, testWraps[1], wraps[1].ToDomain())
	})
}
