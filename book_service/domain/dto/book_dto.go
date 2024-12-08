package dto

import (
	"time"
)

type GetBookListReq struct {
	Query     string `form:"query"`
	QueryBy   string `form:"query_by" default:"any" binding:"oneof=title any"`
	Page      int    `form:"page" default:"1"`
	Limit     int    `form:"limit" default:"10"`
	SortOrder string `form:"sort_order" binding:"required,oneof=asc desc"`
	SortBy    string `form:"sort_by" default:"created_at"`
}

type BookRepo_GetListParams struct {
	AuthorUUID string // use string "null" to query null field
	Query      string
	QueryBy    string
	Page       int
	Limit      int
	SortOrder  string
	SortBy     string
}

type CreateBookReq struct {
	CategoryUUID *string `json:"category_uuid"`
	Title        string  `json:"title" binding:"required"`
	Stock        int64   `json:"stock" binding:"required"`
}

type CreateBookResp struct {
	UUID         string    `json:"uuid"`
	AuthorUUID   string    `json:"author_uuid"`
	CategoryUUID *string   `json:"category_uuid"`
	Title        string    `json:"title"`
	Stock        int64     `json:"stock"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// use "no value" to set nullable field to nil
type PatchBookReq struct {
	CategoryUUID *string `json:"category_uuid"`
	Title        *string `json:"title"`
	Stock        *int64  `json:"stock"`
}

type PatchBookRespData struct {
	UUID         string    `json:"uuid"`
	AuthorUUID   string    `json:"author_uuid"`
	CategoryUUID *string   `json:"category_uuid"`
	Title        string    `json:"title"`
	Stock        int64     `json:"stock"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type DeleteBookRespData struct {
	UUID         string    `json:"uuid"`
	AuthorUUID   string    `json:"author_uuid"`
	CategoryUUID *string   `json:"category_uuid"`
	Title        string    `json:"title"`
	Stock        int64     `json:"stock"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type BulkGetBookTotalByAuthorUUIDsReq struct {
	AuthorUUIDs []string `json:"author_uuids" binding:"required"`
}

type BulkGetBookTotalByAuthorUUIDsRespDataItem struct {
	AuthorUUID string `json:"author_uuid"`
	Total      int64  `json:"total"`
}
