package model

type RegisterUserRequest struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Phone           string `json:"phone"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Bio             string `json:"bio"`
	AvatarUrl       string `json:"avatar_url"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
