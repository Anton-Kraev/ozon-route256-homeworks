package wrap

import (
	errsdomain "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/order"
)

type Wrapper interface {
	WrapOrder(order *order.Order) error
}

func GetWrapper(wrap order.Wrap) (Wrapper, error) {
	switch wrap {
	case order.Package:
		return Package{}, nil
	case order.Box:
		return Box{}, nil
	case order.Tape:
		return Tape{}, nil
	case order.Nowrap:
		return NoWrap{}, nil
	default:
		return nil, errsdomain.ErrUnknownOrderWrapType
	}
}
