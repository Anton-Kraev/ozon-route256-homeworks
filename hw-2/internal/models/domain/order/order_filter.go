package order

import (
	"math"
)

type Filter struct {
	OrdersID     []uint64
	ClientsID    []uint64
	Statuses     []Status
	PageN        uint
	PerPage      uint
	SortedByDate bool
}

func (f *Filter) Init() {
	if f.PerPage == 0 {
		f.PerPage = math.MaxUint
	}
}
