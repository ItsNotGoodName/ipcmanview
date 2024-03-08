package rpcserver

import (
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
)

// ---------- Sort

func decodeSort(v *rpc.Sort) _Sort {
	if v == nil {
		return _Sort{}
	}
	return _Sort{
		Field: v.Field,
		Order: v.Order,
	}
}

type _Sort struct {
	Field string
	Order rpc.Order
}

func (s _Sort) defaultOrder(order rpc.Order) _Sort {
	if s.Order == rpc.Order_ORDER_UNSPECIFIED {
		s.Order = order
	}
	return s
}

func (s _Sort) encode() *rpc.Sort {
	return &rpc.Sort{
		Field: s.Field,
		Order: s.Order,
	}
}

func encodeMonthID(month time.Time) string {
	month = month.UTC()
	return fmt.Sprintf("%02d-%02d", month.Year(), month.Month())
}

func decodeMonthID(month string) time.Time {
	t, err := time.ParseInLocation("2006-01", month, time.UTC)
	if err != nil {
		return time.Time{}
	}
	return t
}

// ---------- Order

func decodeOrderSQL(sql string, o rpc.Order) string {
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

func decodePagePagination(v *rpc.PagePagination) pagination.Page {
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

func encodePagePaginationResult(v pagination.PageResult) *rpc.PagePaginationResult {
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
