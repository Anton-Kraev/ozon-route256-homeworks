package cli

import (
	"errors"
	"flag"
	"fmt"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/requests"
	"strings"
)

func (c *CLI) clientOrders(args []string) error {
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

	req := requests.ClientOrdersRequest{
		ClientID:  clientID,
		LastN:     lastN,
		InStorage: inStorage,
	}
	c.cmdManager.AddTask(func() (string, error) {
		return clientOrdersTask(c, req)
	})

	return nil
}

func clientOrdersTask(cli *CLI, req requests.ClientOrdersRequest) (string, error) {
	orders, err := cli.Service.ClientOrders(req)
	if err != nil {
		return "", err
	}

	result := strings.Builder{}

	result.WriteString("\nOrders list:")
	for _, order := range orders {
		result.WriteString(fmt.Sprintf(
			"  orderID=%d clientID=%d status=%s\n",
			order.OrderID, order.ClientID, order.Status,
		))
	}

	return result.String(), nil
}
