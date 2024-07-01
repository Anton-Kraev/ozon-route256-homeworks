package wrap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
	wrapRepo "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/repository/wrap"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/tests/pg"
)

func TestGetWrapByName(t *testing.T) {
	var (
		ctx     = context.Background()
		names   = []string{"box", "package"}
		badName = "unknown"

		testWraps = []wrap.Wrap{
			{
				Name:   names[0],
				Weight: 1000,
				Cost:   1000,
			},
			{
				Name:   names[1],
				Weight: 2000,
				Cost:   2000,
			},
		}
	)

	pg.DB.SetUp(t, "orders", "wrap")
	defer pg.DB.TearDown(t)
	pg.DB.FillWraps(testWraps)
	repo := wrapRepo.NewWrapRepository(pg.DB.ConnPool)

	t.Run("get_wrap", func(t *testing.T) {
		t.Parallel()

		resp, err := repo.GetWrapByName(ctx, names[0])

		require.NoError(t, err)
		require.NotNil(t, resp)
		AssertEqualWraps(t, testWraps[0], *resp)
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		resp, err := repo.GetWrapByName(ctx, badName)

		require.NoError(t, err)
		assert.Nil(t, resp)
	})
}
