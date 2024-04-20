package common

import (
	"math"
)

type PageLimit struct {
	Page        uint32 `form:"page"`
	Limit       uint16 `form:"limit" `
	OrderBy     string `form:"order_by"`
	Sort        string `form:"sort"  binding:"omitempty,oneof=asc desc"`
	SearchField string `form:"search_field"`
	SearchValue string `form:"search_value"`
}

type Pagination[T any] struct {
	Limit       uint16 `json:"limit"`
	CurrentPage uint32 `json:"current_page"`
	TotalPage   uint32 `json:"total_page"`
	Results     T      `json:"results"`
	TotalItem   uint64 `json:"total_item"`
}

func GetDefaultPageLimit(page uint32, limit uint16) PageLimit {
	return PageLimit{
		Page:        page,
		Limit:       limit,
		OrderBy:     "create_date",
		Sort:        "asc",
		SearchField: "",
		SearchValue: "",
	}
}

// For Combind Pagination and Data into struct.
func CombindPaginationWithData[T any](pageLimit PageLimit, total uint64, results T) *Pagination[T] {
	TotalPage := uint32(math.Ceil(float64(total) / float64(pageLimit.Limit)))
	return &Pagination[T]{
		Limit:       pageLimit.Limit,
		CurrentPage: pageLimit.Page,
		TotalPage:   TotalPage,
		Results:     results,
		TotalItem:   total,
	}
}
