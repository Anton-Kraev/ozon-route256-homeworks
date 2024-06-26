package wrap

import (
	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/errors"
	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
)

type Package struct{}

func (w Package) WrapOrder(order *models.Order) error {
	if order.Weight >= 10000 {
		return errsdomain.ErrOrderWeightExceedsLimit
	}

	order.Cost += 5

	return nil
}
