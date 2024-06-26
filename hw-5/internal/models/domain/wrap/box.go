package wrap

import (
	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-4/internal/models/domain/errors"
	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-4/internal/models/domain/order"
)

type Box struct{}

func (w Box) WrapOrder(order *models.Order) error {
	if order.Weight >= 30000 {
		return errsdomain.ErrOrderWeightExceedsLimit
	}

	order.Cost += 20

	return nil
}
