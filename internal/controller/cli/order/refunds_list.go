package order

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	models "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/order"
)

func (c OrderController) RefundsList(ctx context.Context, args []string) (string, error) {
	var pageN, perPage uint

	fs := flag.NewFlagSet("rlist", flag.ContinueOnError)
	fs.UintVar(&pageN, "pageN", 0, "use --pageN=3")
	fs.UintVar(&perPage, "perPage", 0, "use --perPage=10")

	if err := fs.Parse(args); err != nil {
		return "", fmt.Errorf("%w: %w", errParseArgs, err)
	}

	refunds, err := c.service.RefundsList(ctx, pageN, perPage)
	if err != nil {
		log.Println(err.Error())

		return "", errRefundsList
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
