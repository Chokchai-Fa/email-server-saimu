package dbpostgres

// DBPG struct
type DBPG struct {
}

// NewDBPostgres is the function allocate
func NewDBPostgres() *DBPG {
	return &DBPG{}
}

type SearchCondtionValue struct {
	Mode      string
	Condition string
	Action    string
	Field     string
	Value     interface{}
}

type Pagination struct {
	Limit       uint16 `json:"limit"`
	CurrentPage uint32 `json:"current_page"`
	TotalPage   uint32 `json:"total_page"`
	Results     []any  `json:"results"`
	TotalItem   uint64 `json:"total_item"`
}

type PageLimit struct {
	Page        uint32 `form:"page"`
	Limit       uint16 `form:"limit" `
	OrderBy     string `form:"order_by"`
	Sort        string `form:"sort"  binding:"omitempty,oneof=asc desc"`
	SearchField string `form:"search_field"`
	SearchValue string `form:"search_value"`
}
