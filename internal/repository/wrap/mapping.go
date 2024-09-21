package wrap

import (
	"database/sql"
	"time"

	"gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/wrap"
)

type WrapSchema struct {
	Name      string       `db:"name"`
	MaxWeight uint         `db:"max_weight"`
	Cost      uint         `db:"cost"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

func (r WrapSchema) ToDomain() wrap.Wrap {
	return wrap.Wrap{
		Name:      r.Name,
		MaxWeight: r.MaxWeight,
		Cost:      r.Cost,
	}
}
