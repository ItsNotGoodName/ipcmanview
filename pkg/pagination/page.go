package pagination

import (
	"math"
)

type Page struct {
	Page    int
	PerPage int
}

func (p Page) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func (p Page) Limit() int {
	return p.PerPage
}

type PageResult struct {
	Page       int
	PerPage    int
	TotalPages int
	TotalItems int
}

func (p Page) Result(totalItems int) PageResult {
	totalPage := int(math.Ceil(float64(float64(totalItems) / float64(p.PerPage))))
	if totalPage == 0 {
		totalPage = 1
	}
	page := p.Page
	if page > totalPage {
		page = totalPage + 1
	}
	return PageResult{
		Page:       page,
		PerPage:    p.PerPage,
		TotalPages: totalPage,
		TotalItems: totalItems,
	}
}

func (p PageResult) Overflow() bool {
	return p.Page > p.TotalPages
}

func (p PageResult) HasNext() bool {
	return p.Page < p.TotalPages
}

func (p PageResult) Next() int {
	if !p.HasNext() {
		return p.Page
	}
	return p.Page + 1
}

func (p PageResult) HasPrevious() bool {
	return p.Page > 1
}

func (p PageResult) Previous() int {
	if !p.HasPrevious() {
		return p.Page
	}
	return p.Page - 1
}

func (p PageResult) Seen() int {
	seen := p.Page * p.PerPage
	if seen > p.TotalItems {
		return p.TotalItems
	}
	return seen
}
