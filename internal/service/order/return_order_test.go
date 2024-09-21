package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	errsdomain "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/errors"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/order"
)

func TestOrderService_ReturnOrder(t *testing.T) {
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
		name    string
		fields  fields
		mockFn  func(f fields)
		orderID uint64
		wantErr error
	}{
		{
			name: "err_not_found",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(0)).Return(nil, nil)
			},
			orderID: 0,
			wantErr: errsdomain.ErrOrderNotFound,
		},
		{
			name: "err_not_expired_yet",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(1)).Return(
					&order.Order{StoredUntil: afterNow}, nil,
				)
			},
			orderID: 1,
			wantErr: errsdomain.ErrRetentionPeriodNotExpiredYet,
		},
		{
			name: "err_already_returned",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(2)).Return(
					&order.Order{Status: order.Returned, StoredUntil: beforeNow}, nil,
				)
			},
			orderID: 2,
			wantErr: errsdomain.ErrOrderAlreadyReturned,
		},
		{
			name: "err_delivered",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(3)).Return(
					&order.Order{Status: order.Delivered, StoredUntil: beforeNow}, nil,
				)
			},
			orderID: 3,
			wantErr: errsdomain.ErrOrderDelivered,
		},
		{
			name: "successful_return",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(4)).Return(
					&order.Order{Status: order.Received, StatusChanged: beforeNow, StoredUntil: beforeNow}, nil,
				)
				f.orderRepo.EXPECT().ChangeOrders(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, changes []order.Order) error {
						if changes[0].StatusChanged.Sub(now) < time.Hour &&
							changes[0].Hash == "hash" &&
							changes[0].Status == order.Returned {
							return nil
						}

						return errors.New("")
					})
				f.hashes.EXPECT().GetHash().Return("hash")
			},
			orderID: 4,
			wantErr: nil,
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
			err := service.ReturnOrder(ctx, tt.orderID)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
