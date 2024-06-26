package order

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/helpers"
	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/errors"
)

func (c OrderController) DeliverOrders(ctx context.Context, args []string) (string, error) {
	var ordersStr string

	fs := flag.NewFlagSet("deliver", flag.ContinueOnError)
	fs.StringVar(&ordersStr, "orders", "", "use --orders=1,2,3,4,5")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	orders, err := helpers.StrToUint64Arr(ordersStr)
	if err != nil {
		return "", fmt.Errorf("invalid input format, must be --orders=<id1>,<id2>,<id3>: %v", err)
	}

	for _, order := range orders {
		if order == 0 {
			return "", errors.New("orderID must be positive number")
		}
	}

	err = c.service.DeliverOrders(ctx, orders)
	if err != nil {
		switch {
		case errors.Is(err, errsdomain.ErrOrderNotFound):
			return "", err
		case errors.Is(err, errsdomain.ErrDifferentClients):
			return "", err
		case errors.Is(err, errsdomain.ErrUnexpectedOrderStatus):
			return "", err
		case errors.Is(err, errsdomain.ErrRetentionPeriodExpired):
			return "", err
		default:
			log.Println(err.Error())

			return "", errors.New("can't deliver orders")
		}
	}

	return "", nil
}
