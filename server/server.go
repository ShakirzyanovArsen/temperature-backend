package server

import (
	"net/http"
	"temperature-backend/handler"
	"temperature-backend/middleware"
	"temperature-backend/repository"
	"temperature-backend/service"
)

func Setup() {
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(&userRepo)
	userHandler := handler.NewUserHandler(&userService)
	http.Handle("/user/register", middleware.Post(http.HandlerFunc(userHandler.RegisterUser)))
}