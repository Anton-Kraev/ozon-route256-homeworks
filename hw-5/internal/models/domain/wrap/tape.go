package wrap

import (
	models "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
)

type Tape struct{}

func (w Tape) WrapOrder(order *models.Order) error {
	order.Cost += 1

	return nil
}
