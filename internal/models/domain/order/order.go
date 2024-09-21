package order

import (
	"time"

	errsdomain "gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/errors"
	"gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/wrap"
)

type Order struct {
	OrderID       uint64
	ClientID      uint64
	Weight        uint
	Cost          uint
	WrapType      string
	StoredUntil   time.Time
	Status        Status
	StatusChanged time.Time
	Hash          string
}

func (o *Order) Wrap(wrap wrap.Wrap) error {
	if o.Weight > wrap.MaxWeight {
		return errsdomain.ErrOrderWeightExceedsLimit
	}

	o.WrapType = wrap.Name
	o.Cost += wrap.Cost

	return nil
}

func (o *Order) SetStatus(status Status, timeChanged time.Time) {
	o.Status = status
	o.StatusChanged = timeChanged
}

func (o *Order) SetHash(hash string) {
	o.Hash = hash
}
