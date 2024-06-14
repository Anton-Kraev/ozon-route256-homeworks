package order

import (
	"errors"
	"flag"
	"fmt"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/helpers"
	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/errors"
)

func (c *CLI) deliverOrders(args []string) (string, error) {
	var ordersStr string

	fs := flag.NewFlagSet(deliverOrders, flag.ContinueOnError)
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

	errDeliver := c.Service.DeliverOrders(orders)
	if errDeliver != nil {
		switch {
		case errors.Is(errDeliver, errsdomain.ErrOrderNotFound) ||
			errors.Is(errDeliver, errsdomain.ErrDifferentClients) ||
			errors.Is(errDeliver, errsdomain.ErrUnexpectedOrderStatus) ||
			errors.Is(errDeliver, errsdomain.ErrRetentionPeriodExpired):
			return "", errDeliver
		default:
			return "", errors.New("can't deliver orders")
		}
	}

	return "", nil
}
