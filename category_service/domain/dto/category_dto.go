package dto

import "time"

type CategoryRepo_GetListParams struct {
	Query     string
	QueryBy   string // leave empty to query by any queriable fields
	Offset    int
	Limit     int
	SortOrder string
	SortBy    string
}

type CreateCategoryReq struct {
	Name string `json:"name" binding:"required"`
}

type CreateCategoryRespData struct {
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	Name      string    `json:"name"`
}

type PatchCategoryReq struct {
	Name *string `json:"name,omitempty"`
}

type PatchCategoryRespData struct {
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	Name      string    `json:"name"`
}

type DeleteCategoryRespData struct {
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	Name      string    `json:"name"`
}

type GetCategoryDetailRespData struct {
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	Name      string    `json:"name"`
	BookTotal int64     `json:"book_total"`
}

type GetCategoryListReq struct {
	Query     string `form:"query" default:""`
	QueryBy   string `form:"query_by" default:"any" binding:"oneof=name any"`
	Page      int    `form:"page" default:"1"`
	Limit     int    `form:"limit" default:"10"`
	SortOrder string `form:"sort_order" default:"desc" binding:"required,oneof=asc desc"`
	SortBy    string `form:"sort_by" default:"created_at" binding:"oneof=created_at updated_at name"`
}

type GetListCategoryRespDataItem struct {
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `json:"created_by"`
	Name      string    `json:"name"`
	BookTotal int64     `json:"book_total"`
}

type GetListCategoryRespData struct {
	BasePaginatedData
	Data []GetListCategoryRespDataItem `json:"data"`
}
