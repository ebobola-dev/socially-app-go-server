package model

type UserCreateRequest struct {
	Email string `json:"email"`
}

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}
