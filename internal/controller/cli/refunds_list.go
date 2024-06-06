package cli

import (
	"flag"
	"fmt"

	"gitlab.ozon.dev/antonkraeww/homeworks/internal/domain/requests"
)

func (c CLI) refundsList(args []string) error {
	var pageN, perPage uint

	fs := flag.NewFlagSet(refundsList, flag.ContinueOnError)
	fs.UintVar(&pageN, "pageN", 0, "use --pageN=3")
	fs.UintVar(&perPage, "perPage", 0, "use --perPage=10")

	if err := fs.Parse(args); err != nil {
		return err
	}

	refunds, err := c.Service.RefundsList(requests.RefundsListRequest{PageN: pageN, PerPage: perPage})
	if err == nil {
		fmt.Println("\nRefunds list:")
		for _, refund := range refunds {
			fmt.Printf(
				"  orderID=%d clientID=%d refunded=%s\n",
				refund.OrderID, refund.ClientID, refund.StatusChanged,
			)
		}
	}

	return err
}
