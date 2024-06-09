package cli

import (
	"flag"
	"fmt"
	"gitlab.ozon.dev/antonkraeww/homeworks/internal/models/requests"
	"strings"
)

func (c *CLI) refundsList(args []string) error {
	var pageN, perPage uint

	fs := flag.NewFlagSet(refundsList, flag.ContinueOnError)
	fs.UintVar(&pageN, "pageN", 0, "use --pageN=3")
	fs.UintVar(&perPage, "perPage", 0, "use --perPage=10")

	if err := fs.Parse(args); err != nil {
		return err
	}

	req := requests.RefundsListRequest{PageN: pageN, PerPage: perPage}
	c.cmdManager.AddTask(func() (string, error) {
		return refundsListTask(c, req)
	})

	return nil
}

func refundsListTask(cli *CLI, req requests.RefundsListRequest) (string, error) {
	refunds, err := cli.Service.RefundsList(req)
	if err != nil {
		return "", err
	}

	result := strings.Builder{}

	result.WriteString("\nRefunds list:")
	for _, refund := range refunds {
		result.WriteString(fmt.Sprintf(
			"  orderID=%d clientID=%d refunded=%s\n",
			refund.OrderID, refund.ClientID, refund.StatusChanged,
		))
	}

	return result.String(), nil
}
