package cli

import (
	"errors"
	"flag"
	"fmt"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/requests"
)

func (c CLI) clientOrders(args []string) error {
	var (
		clientID  uint64
		lastN     uint
		inStorage bool
	)

	fs := flag.NewFlagSet(clientOrders, flag.ContinueOnError)
	fs.Uint64Var(&clientID, "clientID", 0, "use --clientID=12345")
	fs.UintVar(&lastN, "lastN", 0, "use --lastN=10")
	fs.BoolVar(&inStorage, "inStorage", false, "use --inStorage")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if clientID == 0 {
		return errors.New("clientID must be positive number")
	}

	orders, err := c.Service.ClientOrders(requests.ClientOrdersRequest{
		ClientID:  clientID,
		LastN:     lastN,
		InStorage: inStorage,
	})
	if err == nil {
		fmt.Println("\nOrders list:")
		for _, order := range orders {
			fmt.Printf(
				"  orderID=%d clientID=%d status=%s\n",
				order.OrderID, order.ClientID, order.Status,
			)
		}
	}

	return err
}
