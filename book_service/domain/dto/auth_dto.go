package dto

type CurrentUser struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
