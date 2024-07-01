package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
)

func TestOrderService_DeliverOrders(t *testing.T) {
	t.Parallel()

	type fields struct {
		orderRepo *MockorderRepository
		hashes    *MockhashGenerator
	}

	var (
		ctx       = context.Background()
		now       = time.Now().UTC()
		afterNow  = now.Add(time.Hour)
		beforeNow = now.Add(-time.Hour)
	)

	tests := []struct {
		name     string
		fields   fields
		mockFn   func(f fields)
		ordersID []uint64
		wantErr  error
	}{
		{
			name: "err_not_found",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrdersByIDs(gomock.Any(), []uint64{0, 1}).Return([]order.Order{{}}, nil)
			},
			ordersID: []uint64{0, 1},
			wantErr:  errsdomain.ErrOrderNotFound,
		},
		{
			name: "err_retention_expired",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrdersByIDs(gomock.Any(), []uint64{2}).Return(
					[]order.Order{{StoredUntil: beforeNow}}, nil,
				)
			},
			ordersID: []uint64{2},
			wantErr:  errsdomain.ErrRetentionPeriodExpired,
		},
		{
			name: "err_not_received",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrdersByIDs(gomock.Any(), []uint64{3}).Return(
					[]order.Order{{StoredUntil: afterNow}}, nil,
				)
			},
			ordersID: []uint64{3},
			wantErr:  errsdomain.ErrUnexpectedOrderStatus,
		},
		{
			name: "err_different_clients",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrdersByIDs(gomock.Any(), []uint64{4, 5}).Return([]order.Order{
					{ClientID: 4, Status: order.Received, StoredUntil: afterNow},
					{ClientID: 5, Status: order.Received, StoredUntil: afterNow},
				}, nil)
				f.hashes.EXPECT().GetHash().Return("hash")
			},
			ordersID: []uint64{4, 5},
			wantErr:  errsdomain.ErrDifferentClients,
		},
		{
			name: "successful_deliver",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrdersByIDs(gomock.Any(), []uint64{6}).Return(
					[]order.Order{{Status: order.Received, StatusChanged: now, StoredUntil: afterNow}}, nil,
				)
				f.orderRepo.EXPECT().ChangeOrders(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, changes []order.Order) error {
						if changes[0].StatusChanged.Sub(now) < time.Hour &&
							changes[0].Hash == "hash" &&
							changes[0].Status == order.Delivered {
							return nil
						}

						return errors.New("")
					})
				f.hashes.EXPECT().GetHash().Return("hash")
			},
			ordersID: []uint64{6},
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tt.fields = fields{
				NewMockorderRepository(ctrl),
				NewMockhashGenerator(ctrl),
			}
			tt.mockFn(tt.fields)

			service := &OrderService{
				orderRepo: tt.fields.orderRepo,
				hashes:    tt.fields.hashes,
			}
			err := service.DeliverOrders(ctx, tt.ordersID)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
