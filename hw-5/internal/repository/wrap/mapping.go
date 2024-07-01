package wrap

import "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"

type WrapSchema struct {
	Name   string `db:"name"`
	Weight uint   `db:"weight"`
	Cost   uint   `db:"cost"`
}

func (r WrapSchema) ToDomain() wrap.Wrap {
	return wrap.Wrap{
		Name:   r.Name,
		Weight: r.Weight,
		Cost:   r.Cost,
	}
}
