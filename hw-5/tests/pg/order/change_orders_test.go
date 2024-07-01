package order

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/middlewares"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
	orderRepo "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/repository/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/tests/pg"
)

func TestChangeOrders(t *testing.T) {
	var (
		now = time.Now().UTC()

		testWraps = []wrap.Wrap{
			{
				Name:   "box",
				Weight: 10,
				Cost:   10,
			},
			{
				Name:   "package",
				Weight: 10,
				Cost:   10,
			},
		}

		testOrders = []order.Order{
			{
				OrderID:       1,
				ClientID:      1,
				Weight:        1,
				Cost:          1,
				Status:        order.Received,
				WrapType:      "box",
				StoredUntil:   now,
				StatusChanged: now,
				Hash:          "hash",
			},
			{
				OrderID:       2,
				ClientID:      2,
				Weight:        2,
				Cost:          2,
				Status:        order.Delivered,
				WrapType:      "package",
				StoredUntil:   now,
				StatusChanged: now,
				Hash:          "hash",
			},
		}

		testOrdersChanged = []order.Order{
			{
				OrderID:       1,
				ClientID:      1,
				Weight:        1,
				Cost:          1,
				Status:        order.Delivered,
				WrapType:      "box",
				StoredUntil:   now,
				StatusChanged: now.Add(time.Hour),
				Hash:          "hash",
			},
			{
				OrderID:       2,
				ClientID:      2,
				Weight:        2,
				Cost:          2,
				Status:        order.Refunded,
				WrapType:      "package",
				StoredUntil:   now,
				StatusChanged: now,
				Hash:          "newhash",
			},
			{
				OrderID:       3,
				ClientID:      3,
				Weight:        3,
				Cost:          3,
				Status:        order.Received,
				WrapType:      "box",
				StoredUntil:   now,
				StatusChanged: now,
				Hash:          "hash",
			},
		}
	)

	pg.DB.SetUp(t, "orders", "wrap")
	defer pg.DB.TearDown(t)
	pg.DB.FillWraps(testWraps)
	pg.DB.FillOrders(testOrders)
	repo := orderRepo.NewOrderRepository(pg.DB.ConnPool)
	txMw := middlewares.NewTransactionMiddleware(pg.DB.ConnPool)

	t.Run("no_changes", func(t *testing.T) {
		ctx := context.Background()
		res, err := txMw.CreateTransactionContext(
			ctx,
			pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead},
			[]string{},
			func(ctx context.Context, args []string) (string, error) {
				return "", repo.ChangeOrders(ctx, testOrders)
			},
		)

		require.NoError(t, err)
		require.Empty(t, res)

		records := pg.DB.GetAllOrders()
		require.NotEmpty(t, records)
		AssertEqualOrders(t, testOrders[0], records[0].ToDomain())
		AssertEqualOrders(t, testOrders[1], records[1].ToDomain())
	})

	t.Run("successful_changes", func(t *testing.T) {
		ctx := context.Background()
		res, err := txMw.CreateTransactionContext(
			ctx,
			pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead},
			[]string{},
			func(ctx context.Context, args []string) (string, error) {
				return "", repo.ChangeOrders(ctx, testOrdersChanged[:2])
			},
		)

		require.NoError(t, err)
		require.Empty(t, res)

		records := pg.DB.GetAllOrders()
		require.NotEmpty(t, records)
		AssertEqualOrders(t, testOrdersChanged[0], records[0].ToDomain())
		AssertEqualOrders(t, testOrdersChanged[1], records[1].ToDomain())
	})

	t.Run("not_found", func(t *testing.T) {
		ctx := context.Background()
		res, err := txMw.CreateTransactionContext(
			ctx,
			pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead},
			[]string{},
			func(ctx context.Context, args []string) (string, error) {
				return "", repo.ChangeOrders(ctx, testOrdersChanged)
			},
		)

		require.NoError(t, err)
		require.Empty(t, res)

		records := pg.DB.GetAllOrders()
		require.NotEmpty(t, records)
		AssertEqualOrders(t, testOrdersChanged[0], records[0].ToDomain())
		AssertEqualOrders(t, testOrdersChanged[1], records[1].ToDomain())
	})
}
