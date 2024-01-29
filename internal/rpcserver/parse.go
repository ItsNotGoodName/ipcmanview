package rpcserver

import (
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
)

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
	if v.PerPage < 1 || v.PerPage > 100 {
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
