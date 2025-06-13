package router

import (
	"net/http"
	"socially-app/internal/config"
	"socially-app/internal/handler"
	logger "socially-app/internal/util"

	"github.com/gorilla/mux"
)

func New(cfg *config.Config, log logger.ILogger) *http.Server {
	r := mux.NewRouter()

	userHandler := handler.NewUserHandler(log)

	r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}
}
