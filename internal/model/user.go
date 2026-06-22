package model

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username"`
	PasswordHash []byte `json:"-"`
	Salt         []byte `json:"-"`
	AvatarUrl    string `json:"avatar_url,omitempty"`
	Bio          string `json:"bio,omitempty"`
}

type UserIdentifier struct {
	Username string
	Email    string
	Phone    string
}
