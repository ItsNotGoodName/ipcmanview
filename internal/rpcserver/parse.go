package rpcserver

import (
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
)

// ---------- Sort

type rpcSort struct {
	Field string
	Order rpc.Order
}

func (s rpcSort) withDefaultOrder(order rpc.Order) rpcSort {
	if s.Order == rpc.Order_ORDER_UNSPECIFIED {
		s.Order = order
	}
	return s
}

func parseSort(v *rpc.Sort) rpcSort {
	if v == nil {
		return rpcSort{}
	}
	return rpcSort{
		Field: v.Field,
		Order: v.Order,
	}
}

// ---------- Order

func parseOrderSQL(sql string, o rpc.Order) string {
	switch o {
	case rpc.Order_DESC:
		return sql + " DESC"
	case rpc.Order_ASC:
		return sql + " ASC"
	default:
		return sql
	}
}

// ---------- Page

func parsePagePagination(v *rpc.PagePagination) pagination.Page {
	var (
		page    int
		perPage int
	)
	if v != nil {
		page = int(v.Page)
		perPage = int(v.PerPage)
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	return pagination.Page{
		Page:    page,
		PerPage: perPage,
	}
}

func parsePagePaginationResult(v pagination.PageResult) *rpc.PagePaginationResult {
	return &rpc.PagePaginationResult{
		Page:         int32(v.Page),
		PerPage:      int32(v.PerPage),
		TotalPages:   int32(v.TotalPages),
		TotalItems:   int64(v.TotalItems),
		SeenItems:    int64(v.Seen()),
		PreviousPage: int32(v.Previous()),
		NextPage:     int32(v.Next()),
	}
}
