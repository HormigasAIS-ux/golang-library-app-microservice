package dto

import "time"

type GetAuthorListReq struct {
	Query     string `form:"query"`
	QueryBy   string `form:"query_by" default:"any" binding:"oneof=first_name last_name birth_date any"`
	Page      int    `form:"page" default:"1"`
	Limit     int    `form:"limit" default:"10"`
	SortOrder string `form:"sort_order" binding:"required,oneof=asc desc"`
	SortBy    string `form:"sort_by" default:"created_at"`
}

type GetAuthorListRespDataItem struct {
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate *string   `json:"birth_date"`
	Bio       *string   `json:"bio"`
}

type GetAuthorListRespData struct {
	BasePaginatedData
	Data []GetAuthorListRespDataItem `json:"data"`
}

type AuthorRepo_GetListParams struct {
	Query     string
	QueryBy   string
	Page      int
	Limit     int
	SortOrder string
	SortBy    string
}
