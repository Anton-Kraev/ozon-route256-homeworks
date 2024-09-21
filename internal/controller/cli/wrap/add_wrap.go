package wrap

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	errsdomain "github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/errors"
	"github.com/Anton-Kraev/ozon-route256-homeworks/internal/models/domain/wrap"
)

var (
	errParseArgs  = errors.New("can't parse args")
	errAddNewWrap = errors.New("can't add new wrap")
)

func (c WrapController) AddWrap(ctx context.Context, args []string) (string, error) {
	var (
		name            string
		maxWeight, cost uint
	)

	fs := flag.NewFlagSet("addwrap", flag.ContinueOnError)
	fs.StringVar(&name, "name", "", "use --name=box")
	fs.UintVar(&maxWeight, "weight", 0, "use --weight=1000")
	fs.UintVar(&cost, "cost", 0, "use --cost=10")

	if err := fs.Parse(args); err != nil {
		return "", fmt.Errorf("%w: %w", errParseArgs, err)
	}

	if name == "" {
		return "", fmt.Errorf("%w: wrap name is required", errParseArgs)
	}

	err := c.wrapService.AddWrap(ctx, wrap.Wrap{Name: name, MaxWeight: maxWeight, Cost: cost})
	if err != nil {
		switch {
		case errors.Is(err, errsdomain.ErrWrapAlreadyExists):
			return "", err
		default:
			log.Println(err.Error())

			return "", errAddNewWrap
		}
	}

	return "", nil
}
