package wrap

import (
	"context"
	"errors"
	"flag"
	"log"

	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
)

func (c WrapController) AddWrap(ctx context.Context, args []string) (string, error) {
	var (
		name         string
		weight, cost uint
	)

	fs := flag.NewFlagSet("addwrap", flag.ContinueOnError)
	fs.StringVar(&name, "name", "", "use --name=box")
	fs.UintVar(&weight, "weight", 0, "use --weight=1000")
	fs.UintVar(&cost, "cost", 0, "use --cost=10")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if name == "" {
		return "", errors.New("wrap name is required")
	}

	err := c.wrapService.AddWrap(ctx, wrap.Wrap{Name: name, Weight: weight, Cost: cost})
	if err != nil {
		switch {
		case errors.Is(err, errsdomain.ErrWrapAlreadyExists):
			return "", err
		default:
			log.Println(err.Error())

			return "", errors.New("can't add new wrap")
		}
	}

	return "", nil
}
