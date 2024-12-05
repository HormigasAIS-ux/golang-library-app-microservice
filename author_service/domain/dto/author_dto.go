package dto

import (
	"time"

	"github.com/google/uuid"
)

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
	BookTotal int64     `json:"book_total"`
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

type CreateNewAuthorReq struct {
	UserUUID  *string `json:"user_uuid"` // required for create new author by auth service, optional for client
	Email     string  `json:"email" binding:"required,email"`
	Username  string  `json:"username" binding:"required"`
	Password  string  `json:"password" binding:"required"`
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name"`
	BirthDate *string `json:"birth_date"`
	Bio       *string `json:"bio"`
	Role      string  `json:"role" binding:"required,oneof=admin user"`
}

type CreateNewAuthorRespData struct {
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

type EditAuthorReq struct {
	Email     *string `json:"email" binding:"email"`
	Username  *string `json:"username"`
	Password  *string `json:"password"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	BirthDate *string `json:"birth_date"`
	Bio       *string `json:"bio"`
	Role      *string `json:"role" binding:"oneof=admin user"`
}

type EditAuthorRespData struct {
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

type DeleteAuthorRespData struct {
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

type GetAuthorDetailRespData struct {
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

type GetAuthorByUserUUIDRespData struct {
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserUUID  uuid.UUID `json:"user_uuid"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate *string   `json:"birth_date"`
	Bio       *string   `json:"bio"`
}
