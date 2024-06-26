package wrap

import models "gitlab.ozon.dev/antonkraeww/homeworks/hw-4/internal/models/domain/order"

type NoWrap struct{}

func (w NoWrap) WrapOrder(_ *models.Order) error {
	return nil
}
