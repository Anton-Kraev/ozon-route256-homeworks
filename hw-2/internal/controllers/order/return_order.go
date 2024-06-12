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

	errReturn := c.Service.ReturnOrder(orderID)
	if errReturn != nil {
		switch {
		case errors.Is(errReturn, errsdomain.ErrOrderNotFound) ||
			errors.Is(errReturn, errsdomain.ErrRetentionPeriodNotExpiredYet) ||
			errors.Is(errReturn, errsdomain.ErrOrderAlreadyReturned) ||
			errors.Is(errReturn, errsdomain.ErrOrderDelivered):
			return "", errReturn
		default:
			return "", errors.New("can't return order")
		}
	}

	return "", nil
}
