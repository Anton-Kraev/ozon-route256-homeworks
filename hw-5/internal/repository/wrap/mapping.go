package wrap

import (
	"database/sql"
	"time"

	"gitlab.ozon.dev/antonkraeww/homeworks/hw-5/internal/models/domain/wrap"
)

type WrapSchema struct {
	Name      string       `db:"name"`
	Weight    uint         `db:"weight"`
	Cost      uint         `db:"cost"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

func (r WrapSchema) ToDomain() wrap.Wrap {
	return wrap.Wrap{
		Name:   r.Name,
		Weight: r.Weight,
		Cost:   r.Cost,
	}
}
