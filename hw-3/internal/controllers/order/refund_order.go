package order

import (
	"context"
	"errors"
	"flag"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/errors"
)

func (c *CLI) refundOrder(ctx context.Context, args []string) (string, error) {
	var orderID, clientID uint64

	fs := flag.NewFlagSet(refundOrder, flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")
	fs.Uint64Var(&clientID, "clientID", 0, "use --clientID=67890")

	if err := fs.Parse(args); err != nil {
		return "", err
	}
	if orderID == 0 {
		return "", errors.New("orderID must be positive number")
	}
	if clientID == 0 {
		return "", errors.New("clientID must be positive number")
	}

	err := c.Service.RefundOrder(ctx, orderID, clientID)
	if err != nil {
		switch {
		case errors.Is(err, errsdomain.ErrOrderNotFound):
			return "", err
		case errors.Is(err, errsdomain.ErrOrderAlreadyRefunded):
			return "", err
		case errors.Is(err, errsdomain.ErrOrderNotDeliveredYet):
			return "", err
		case errors.Is(err, errsdomain.ErrOrderDeliveredLongAgo):
			return "", err
		default:
			return "", errors.New("can't refund order")
		}
	}

	return "", nil
}
