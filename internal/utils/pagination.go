package utils

import "math"

// 1. INPUT: Hứng dữ liệu từ URL
type PaginationQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

// 2. OUTPUT: Cấu trúc JSON trả về cho Frontend
type PaginationResponse struct {
	Data   interface{} `json:"data"`
	Paging PagingInfo  `json:"paging"`
}

type PagingInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

func (q *PaginationQuery) GetOffset() (int, int) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	offset := (q.Page - 1) * q.Limit
	return offset, q.Limit
}

func NewPaginationResponse(data interface{}, page, limit int, totalItems int64) PaginationResponse {
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	// Fix hiển thị page nếu user gửi page <= 0
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	return PaginationResponse{
		Data: data,
		Paging: PagingInfo{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}
}
