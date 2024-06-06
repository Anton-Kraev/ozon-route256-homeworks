package cli

import (
	"errors"
	"flag"
	"strconv"
	"strings"

	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/requests"
)

func (c CLI) deliverOrders(args []string) error {
	var (
		ordersStr  string
		ordersList []uint64
	)

	fs := flag.NewFlagSet(deliverOrders, flag.ContinueOnError)
	fs.StringVar(&ordersStr, "orders", "", "use --orders=1,2,3,4,5")

	if err := fs.Parse(args); err != nil {
		return err
	}

	for _, orderStr := range strings.Split(ordersStr, ",") {
		orderID, err := strconv.ParseUint(orderStr, 10, 64)

		if err != nil {
			return err
		}
		if orderID == 0 {
			return errors.New("orderID must be positive number")
		}

		ordersList = append(ordersList, orderID)
	}

	return c.Service.DeliverOrders(requests.DeliverOrdersRequest{OrdersID: ordersList})
}
