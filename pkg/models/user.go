package models

type User struct {
	Email string `json:"email"`
}

func NewUser(email string) *User {
	return &User{Email: email}
}
