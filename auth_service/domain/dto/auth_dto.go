package dto

type CurrentUser struct {
	ID       string `json:"id"` // user uuid
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}
