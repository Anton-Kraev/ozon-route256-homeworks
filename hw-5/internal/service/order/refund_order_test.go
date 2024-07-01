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

func TestOrderService_RefundOrder(t *testing.T) {
	t.Parallel()

	type (
		fields struct {
			orderRepo *MockorderRepository
			hashes    *MockhashGenerator
		}

		args struct {
			orderID  uint64
			clientID uint64
		}
	)

	var (
		ctx = context.Background()
		now = time.Now().UTC()
	)

	tests := []struct {
		name    string
		fields  fields
		mockFn  func(f fields)
		args    args
		wantErr error
	}{
		{
			name: "err_not_found",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(0)).Return(nil, nil)
			},
			args:    args{0, 0},
			wantErr: errsdomain.ErrOrderNotFound,
		},
		{
			name: "err_not_found_client",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(1)).Return(
					&order.Order{ClientID: 1}, nil,
				)
			},
			args:    args{1, 0},
			wantErr: errsdomain.ErrOrderNotFound,
		},
		{
			name: "err_already_refunded",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(2)).Return(
					&order.Order{Status: order.Refunded}, nil,
				)
			},
			args:    args{orderID: 2},
			wantErr: errsdomain.ErrOrderAlreadyRefunded,
		},
		{
			name: "err_received",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(3)).Return(
					&order.Order{Status: order.Received}, nil,
				)
			},
			args:    args{orderID: 3},
			wantErr: errsdomain.ErrOrderNotDeliveredYet,
		},
		{
			name: "err_returned",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(4)).Return(
					&order.Order{Status: order.Returned}, nil,
				)
			},
			args:    args{orderID: 4},
			wantErr: errsdomain.ErrOrderNotDeliveredYet,
		},
		{
			name: "err_delivered_long_ago",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(5)).Return(
					&order.Order{Status: order.Delivered, StatusChanged: now.Add(-time.Hour * 1000)}, nil,
				)
			},
			args:    args{orderID: 5},
			wantErr: errsdomain.ErrOrderDeliveredLongAgo,
		},
		{
			name: "successful_refund",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(6)).Return(
					&order.Order{Status: order.Delivered, StatusChanged: now, StoredUntil: now}, nil,
				)
				f.orderRepo.EXPECT().ChangeOrders(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, changes []order.Order) error {
						if changes[0].StatusChanged.Sub(now) < time.Hour &&
							changes[0].Status == order.Refunded &&
							changes[0].Hash == "hash" {
							return nil
						}

						return errors.New("")
					})
				f.hashes.EXPECT().GetHash().Return("hash")
			},
			args:    args{orderID: 6},
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
			err := service.RefundOrder(ctx, tt.args.orderID, tt.args.clientID)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
