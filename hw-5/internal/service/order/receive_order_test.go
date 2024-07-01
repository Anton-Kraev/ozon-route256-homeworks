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
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
)

func TestOrderService_ReceiveOrder(t *testing.T) {
	t.Parallel()

	type (
		fields struct {
			orderRepo *MockorderRepository
			wrapRepo  *MockwrapRepository
			hashes    *MockhashGenerator
		}

		args struct {
			wrapType string
			order    order.Order
		}
	)

	var (
		ctx      = context.Background()
		now      = time.Now().UTC()
		wrapName = "wrap"
	)

	tests := []struct {
		name    string
		fields  fields
		mockFn  func(f fields)
		args    args
		wantErr error
	}{
		{
			name:    "err_retention_in_past",
			mockFn:  func(f fields) {},
			args:    args{"", order.Order{StoredUntil: now.Add(-time.Hour)}},
			wantErr: errsdomain.ErrRetentionTimeInPast,
		},
		{
			name: "err_already_exists",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(1)).Return(
					&order.Order{OrderID: 1}, nil,
				)
			},
			args:    args{"", order.Order{OrderID: 1, StoredUntil: now.Add(time.Hour)}},
			wantErr: errsdomain.ErrOrderIDNotUnique,
		},
		{
			name: "err_wrap_not_found",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(2)).Return(nil, nil)
				f.wrapRepo.EXPECT().GetWrapByName(gomock.Any(), wrapName).Return(nil, nil)
			},
			args:    args{wrapName, order.Order{OrderID: 2, StoredUntil: now.Add(time.Hour)}},
			wantErr: errsdomain.ErrWrapNotFound,
		},
		{
			name: "successful_receive",
			mockFn: func(f fields) {
				f.orderRepo.EXPECT().GetOrderByID(gomock.Any(), uint64(3)).Return(nil, nil)
				f.wrapRepo.EXPECT().GetWrapByName(gomock.Any(), wrapName).Return(
					&wrap.Wrap{Name: wrapName, Weight: 100, Cost: 10}, nil,
				)
				f.orderRepo.EXPECT().AddOrder(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, orderToAdd order.Order) error {
						if orderToAdd.StatusChanged.Sub(now) < time.Hour &&
							orderToAdd.Status == order.Received &&
							orderToAdd.Hash == "hash" &&
							orderToAdd.Weight == 1000 &&
							orderToAdd.Cost == 100 &&
							orderToAdd.WrapType == wrapName {
							return nil
						}

						return errors.New("")
					})
				f.hashes.EXPECT().GetHash().Return("hash")
			},
			args: args{wrapName, order.Order{
				OrderID: 3, Weight: 900, Cost: 90, StoredUntil: now.Add(time.Hour),
			}},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			tt.fields = fields{
				orderRepo: NewMockorderRepository(ctrl),
				wrapRepo:  NewMockwrapRepository(ctrl),
				hashes:    NewMockhashGenerator(ctrl),
			}
			tt.mockFn(tt.fields)

			service := &OrderService{
				orderRepo: tt.fields.orderRepo,
				wrapRepo:  tt.fields.wrapRepo,
				hashes:    tt.fields.hashes,
			}
			err := service.ReceiveOrder(ctx, tt.args.wrapType, tt.args.order)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
