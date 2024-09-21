package wrap

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	errsdomain "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/errors"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/wrap"
)

func TestWrapService_AddWrap(t *testing.T) {
	t.Parallel()

	var (
		ctx     = context.Background()
		oldWrap = "old"
		newWrap = "new"
	)

	t.Run("add_existing", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockRepo := NewMockwrapRepository(ctrl)
		service := NewWrapService(mockRepo)
		mockRepo.EXPECT().GetWrapByName(gomock.Any(), oldWrap).Return(&wrap.Wrap{Name: oldWrap}, nil)

		err := service.AddWrap(ctx, wrap.Wrap{Name: oldWrap})

		require.ErrorIs(t, err, errsdomain.ErrWrapAlreadyExists)
	})

	t.Run("add_new", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockRepo := NewMockwrapRepository(ctrl)
		service := NewWrapService(mockRepo)
		mockRepo.EXPECT().GetWrapByName(gomock.Any(), newWrap).Return(nil, nil)
		mockRepo.EXPECT().AddWrap(gomock.Any(), wrap.Wrap{Name: newWrap}).Return(nil)

		err := service.AddWrap(ctx, wrap.Wrap{Name: newWrap})

		require.NoError(t, err)
	})
}
