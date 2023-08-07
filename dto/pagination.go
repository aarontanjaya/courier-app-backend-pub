package dto

type Pagination struct {
	PageSize   int64 `json:"page_size"`
	TotalCount int64 `json:"total_count"`
	Page       int64 `json:"page"`
	PageCount  int64 `json:"page_count"`
}

type PaginationRequest struct {
	Search   string `json:"search"`
	PageSize int64  `json:"page_size" binding:"numeric"`
	Page     int64  `json:"page" binding:"numeric"`
}

type PaginationSortRequest struct {
	Search        string `json:"search"`
	SortDirection string `json:"sort_dir"`
	SortField     string `json:"sort_field"`
	PageSize      int64  `json:"page_size" binding:"numeric"`
	Page          int64  `json:"page" binding:"numeric"`
}
type AddressPaginationRequest struct {
	UserId   *uint  `json:"user_id"`
	Search   string `json:"search"`
	PageSize int64  `json:"page_size" binding:"numeric"`
	Page     int64  `json:"page" binding:"numeric"`
}
