package dto

type BookBorrowRepo_GetListParams struct {
	AuthorUUID string
	Query      string
	QueryBy    string
	Page       int
	Limit      int
	SortOrder  string
	SortBy     string
}
