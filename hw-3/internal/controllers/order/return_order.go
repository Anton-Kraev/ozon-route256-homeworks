package order

import (
	"errors"
	"flag"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/errors"
)

func (c *CLI) returnOrder(args []string) (string, error) {
	var orderID uint64

	fs := flag.NewFlagSet(returnOrder, flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")

	if err := fs.Parse(args); err != nil {
		return "", err
	}
	if orderID == 0 {
		return "", errors.New("orderID must be positive number")
	}

	err := c.Service.ReturnOrder(orderID)
	if err != nil {
		switch {
		case errors.Is(err, errsdomain.ErrOrderNotFound):
			return "", err
		case errors.Is(err, errsdomain.ErrRetentionPeriodNotExpiredYet):
			return "", err
		case errors.Is(err, errsdomain.ErrOrderAlreadyReturned):
			return "", err
		case errors.Is(err, errsdomain.ErrOrderDelivered):
			return "", err
		default:
			return "", errors.New("can't return order")
		}
	}

	return "", nil
}
