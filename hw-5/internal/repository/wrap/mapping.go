package wrap

import "gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"

type wrapSchema struct {
	Name   string `db:"name"`
	Weight uint   `db:"weight"`
	Cost   uint   `db:"cost"`
}

func (r wrapSchema) toDomain() wrap.Wrap {
	return wrap.Wrap{
		Name:   r.Name,
		Weight: r.Weight,
		Cost:   r.Cost,
	}
}
