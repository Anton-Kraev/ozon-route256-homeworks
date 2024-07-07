package order

import (
	"math"
)

type Filter struct {
	ClientID     uint64
	Statuses     []Status
	PageN        uint
	PerPage      uint
	SortedByDate bool
}

func (f *Filter) Init() {
	if f.PerPage == 0 {
		f.PerPage = math.MaxUint32
	}
}
