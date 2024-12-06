package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetBookListReq struct {
	Query     string `form:"query"`
	QueryBy   string `form:"query_by" default:"any" binding:"oneof=title any"`
	Page      int    `form:"page" default:"1"`
	Limit     int    `form:"limit" default:"10"`
	SortOrder string `form:"sort_order" binding:"required,oneof=asc desc"`
	SortBy    string `form:"sort_by" default:"created_at"`
}

type GetBookListRespDataItem struct {
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate *string   `json:"birth_date"`
	Bio       *string   `json:"bio"`
	BookTotal int64     `json:"book_total"`
}

type GetBookListRespData struct {
	BasePaginatedData
	Data []GetBookListRespDataItem `json:"data"`
}

type BookRepo_GetListParams struct {
	AuthorUUID string
	Query      string
	QueryBy    string
	Page       int
	Limit      int
	SortOrder  string
	SortBy     string
}

type CreateNewBookReq struct {
	UserUUID  *string `json:"user_uuid"` // required for create new Book by auth service, optional for client
	Email     string  `json:"email" binding:"required,email"`
	Username  string  `json:"username" binding:"required"`
	Password  string  `json:"password" binding:"required"`
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name"`
	BirthDate *string `json:"birth_date"`
	Bio       *string `json:"bio"`
	Role      string  `json:"role" binding:"required,oneof=admin user"`
}

type CreateNewBookRespData struct {
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserUUID  uuid.UUID `json:"user_uuid"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate *string   `json:"birth_date"`
	Bio       *string   `json:"bio"`
	Role      string    `json:"role"`
}

type EditBookReq struct {
	Email     *string `json:"email" binding:"email"`
	Username  *string `json:"username"`
	Password  *string `json:"password"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	BirthDate *string `json:"birth_date"`
	Bio       *string `json:"bio"`
	Role      *string `json:"role" binding:"oneof=admin user"`
}

type EditBookRespData struct {
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserUUID  uuid.UUID `json:"user_uuid"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate *string   `json:"birth_date"`
	Bio       *string   `json:"bio"`
	Role      string    `json:"role"`
}

type DeleteBookRespData struct {
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserUUID  uuid.UUID `json:"user_uuid"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate *string   `json:"birth_date"`
	Bio       *string   `json:"bio"`
	Role      string    `json:"role"`
}

type GetBookDetailRespData struct {
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserUUID  uuid.UUID `json:"user_uuid"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate *string   `json:"birth_date"`
	Bio       *string   `json:"bio"`
	Role      string    `json:"role"`
	BookTotal int64     `json:"book_total"`
}

//////////////////////////////

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
