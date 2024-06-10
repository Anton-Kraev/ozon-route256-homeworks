package models

import (
	"math"
)

type OrderFilter struct {
	OrdersID     []uint64
	ClientsID    []uint64
	Statuses     []Status
	PageN        uint
	PerPage      uint
	SortedByDate bool
}

func (f *OrderFilter) Init() {
	if f.OrdersID == nil {
		f.OrdersID = []uint64{}
	}
	if f.ClientsID == nil {
		f.ClientsID = []uint64{}
	}
	if f.Statuses == nil {
		f.Statuses = []Status{}
	}
	if f.PerPage == 0 {
		f.PerPage = math.MaxUint
	}
}
