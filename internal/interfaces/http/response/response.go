package response

import (
	"net/http"
)

type SimpleResponse[T any] struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data"`
}

func NewSimpleResponse[T any](data T) *SimpleResponse[T] {
	return &SimpleResponse[T]{
		Success: true,
		Code:    http.StatusOK,
		Data:    data,
	}
}

func NewErrorResponse(code int, message string) *SimpleResponse[any] {
	return &SimpleResponse[any]{
		Success: false,
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PageInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalCount int `json:"total_count"`
}

type PaginatedResponse[T any] struct {
	Success  bool      `json:"success"`
	Code     int       `json:"code"`
	Message  string    `json:"message,omitempty"`
	Data     []T       `json:"data"`
	PageInfo *PageInfo `json:"page_info"`
}

func NewPaginatedResponse[T any](data []T, pagination *Pagination, count int) *PaginatedResponse[T] {
	if data == nil {
		data = make([]T, 0)
	}

	page := 1
	pageSize := count
	if pagination != nil {
		page = pagination.Page
		pageSize = pagination.PageSize
	}

	return &PaginatedResponse[T]{
		Success: true,
		Code:    http.StatusOK,
		Data:    data,
		PageInfo: &PageInfo{
			Page:       page,
			PageSize:   pageSize,
			TotalCount: count,
		},
	}
}
