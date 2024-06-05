package cli

import (
	"errors"
	"flag"
)

func (c CLI) returnOrder(args []string) error {
	var orderID uint64

	fs := flag.NewFlagSet(returnOrder, flag.ContinueOnError)
	fs.Uint64Var(&orderID, "orderID", 0, "use --orderID=12345")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if orderID == 0 {
		return errors.New("orderID must be positive number")
	}

	return c.Module.ReturnOrder(orderID)
}
