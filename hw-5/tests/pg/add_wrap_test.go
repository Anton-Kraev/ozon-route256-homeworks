package pg

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/middlewares"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
	wrapRepo "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/repository/wrap"
)

func TestAddWrap(t *testing.T) {
	var (
		ctx = context.Background()

		testWraps = []wrap.Wrap{
			{
				Name:   "name3",
				Weight: 1000,
				Cost:   1000,
			},
			{
				Name:   "name4",
				Weight: 2000,
				Cost:   2000,
			},
		}
	)

	DB.SetUp(t, "orders", "wrap")
	defer DB.TearDown(t)
	DB.FillWraps(testWraps[:1])

	repo := wrapRepo.NewWrapRepository(DB.ConnPool)
	txMw := middlewares.NewTransactionMiddleware(DB.ConnPool)

	t.Run("add_existing", func(t *testing.T) {
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

		wraps := DB.GetAllWraps()
		require.Len(t, wraps, 2)
		AssertEqualWraps(t, testWraps[0], wraps[0].ToDomain())
		AssertEqualWraps(t, testWraps[1], wraps[1].ToDomain())
	})
}
