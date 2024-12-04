package dto

import (
	"google.golang.org/grpc/codes"
)

type BaseJSONResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  string      `json:"detail"`
	Data    interface{} `json:"data"`
}

type BaseGrpcResp struct {
	Code    codes.Code  `json:"code"`
	Message string      `json:"message"`
	Detail  string      `json:"detail"`
	Data    interface{} `json:"data"`
}

type BasePaginatedData struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int64 `json:"total_page"`
	TotalData   int64 `json:"total_data"`
}

func (s *BasePaginatedData) Set(
	page int,
	limit int,
	count int64,
) {
	if page == 0 {
		s.CurrentPage = 1
	} else {
		s.CurrentPage = page
	}

	if page != 0 && page != 0 && count > 0 {
		s.TotalPage = int64((count + int64(limit) - 1) / int64(limit))
	}

	s.TotalData = count
}
