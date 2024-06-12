package order

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/order"
)

func (c *CLI) clientOrders(args []string) (string, error) {
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
		return "", err
	}
	if clientID == 0 {
		return "", errors.New("clientID must be positive number")
	}

	orders, err := c.Service.ClientOrders(clientID, lastN, inStorage)
	if err != nil {
		return "", fmt.Errorf("can't get client orders: %v", err)
	}

	return clientOrdersToString(orders), nil
}

func clientOrdersToString(orders []models.Order) string {
	result := strings.Builder{}
	result.WriteString("\nOrders list:\n")

	for _, order := range orders {
		result.WriteString(fmt.Sprintf(
			"  orderID=%d clientID=%d status=%s\n",
			order.OrderID, order.ClientID, order.Status,
		))
	}

	return result.String()
}
