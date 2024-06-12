package order

import (
	"flag"
	"fmt"
	"strings"

	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/domain/order"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-2/internal/models/requests"
)

func (c *CLI) refundsList(args []string) (string, error) {
	var pageN, perPage uint

	fs := flag.NewFlagSet(refundsList, flag.ContinueOnError)
	fs.UintVar(&pageN, "pageN", 0, "use --pageN=3")
	fs.UintVar(&perPage, "perPage", 0, "use --perPage=10")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	refunds, err := c.Service.RefundsList(requests.RefundsListRequest{PageN: pageN, PerPage: perPage})
	if err != nil {
		return "", err
	}

	return refundsListToString(refunds), nil
}

func refundsListToString(refunds []models.Order) string {
	result := strings.Builder{}
	result.WriteString("\nRefunds list:\n")

	for _, refund := range refunds {
		result.WriteString(fmt.Sprintf(
			"  orderID=%d clientID=%d refunded=%s\n",
			refund.OrderID, refund.ClientID, refund.StatusChanged,
		))
	}

	return result.String()
}
