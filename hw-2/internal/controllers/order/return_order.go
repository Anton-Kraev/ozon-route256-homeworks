package order

import (
	"errors"
	"flag"
	"fmt"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/requests"
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

	errReturn := c.Service.ReturnOrder(requests.ReturnOrderRequest{OrderID: orderID})
	if errReturn != nil {
		return "", fmt.Errorf("can't return order: %v", errReturn)
	}

	return "", nil
}
