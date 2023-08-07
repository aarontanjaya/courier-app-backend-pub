package util

import (
	"courier-app/dto"
	"math"
)

var (
	DefaultPageSize int64 = 10
	DefaultPage     int64 = 1
)

func ValidatePaginationRequest(req *dto.PaginationRequest) {
	req.Search = "%" + req.Search + "%"
	if req.PageSize <= 0 {
		req.PageSize = DefaultPageSize
	}

	if req.Page <= 0 {
		req.Page = DefaultPage
	}
}

func BuildPaginationObject(req *dto.PaginationRequest, total int64) *dto.Pagination {
	var pagination *dto.Pagination
	var currentPageSize int64
	var maxPage int64

	currentPageSize = req.PageSize
	if req.PageSize <= 0 {
		currentPageSize = DefaultPageSize
	}

	maxPage = int64(math.Ceil(float64(total) / float64(currentPageSize)))
	pagination = &dto.Pagination{
		PageSize:   currentPageSize,
		TotalCount: total,
		Page:       req.Page,
		PageCount:  maxPage,
	}
	return pagination
}
