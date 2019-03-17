package server

import (
	"net/http"
	"temperature-backend/handler"
	m "temperature-backend/middleware"
	"temperature-backend/repository"
	"temperature-backend/service"
)

func Setup() *http.ServeMux {
	userRepo := repository.NewUserRepository()
	deviceRepo := repository.NewDeviceRepository()

	userService := service.NewUserService(&userRepo)
	deviceService := service.NewDeviceService(&userRepo, &deviceRepo)

	userHandler := handler.NewUserHandler(&userService)
	deviceHandler := handler.NewDeviceHandler(&deviceService)

	mux := http.NewServeMux()
	mux.Handle("/user/register", m.Post(http.HandlerFunc(userHandler.RegisterUser)))
	mux.Handle("/device/register", m.Post(http.HandlerFunc(deviceHandler.RegisterDevice)))
	return mux
}
