package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"socially-app/internal/model"
	logger "socially-app/internal/util"
)

type UserHandler struct {
	log logger.ILogger
}

func NewUserHandler(log logger.ILogger) IUserHandler {
	return &UserHandler{log: log}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	h.log.Error(errors.New("test error"))
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`[{"id":1,"username":"alice"}]`))
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	var req model.UserCreateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	h.log.Info("Creating user: " + req.Email)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"user created"}`))
}
