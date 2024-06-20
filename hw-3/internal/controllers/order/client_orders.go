package order

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-3/internal/models/domain/order"
)

func (c *CLI) clientOrders(ctx context.Context, args []string) (string, error) {
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

	orders, err := c.Service.ClientOrders(ctx, clientID, lastN, inStorage)
	if err != nil {
		log.Println(err.Error())

		return "", errors.New("can't get client orders")
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
