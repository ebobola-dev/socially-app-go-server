package handler

import "net/http"

type IUserHandler interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
}
