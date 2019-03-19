package server

import (
	"net/http"
	"temperature-backend/handler"
	m "temperature-backend/middleware"
	"temperature-backend/repository"
	"temperature-backend/service"
)

func Setup() *http.ServeMux {
	userHandler, deviceHandler, deviceDataHandler := initHandlers()

	mux := http.NewServeMux()
	mux.Handle("/user/register", m.Post(http.HandlerFunc(userHandler.RegisterUser)))

	mux.Handle("/device/", m.Get(m.Auth(http.HandlerFunc(deviceHandler.GetDataList))))
	mux.Handle("/device/register", m.Post(http.HandlerFunc(deviceHandler.RegisterDevice)))
	mux.Handle("/device/list", m.Get(m.Auth(http.HandlerFunc(deviceHandler.GetDeviceList))))

	mux.Handle("/device/data", m.Post(m.Auth(http.HandlerFunc(deviceDataHandler.PushData))))
	return mux
}

func initHandlers() (handler.UserHandler, handler.DeviceHandler, handler.DeviceDataHandler) {
	userRepo := repository.NewUserRepository()
	deviceRepo := repository.NewDeviceRepository()
	deviceDataRepo := repository.NewDeviceDataRepository()

	userService := service.NewUserService(&userRepo)
	deviceService := service.NewDeviceService(&userRepo, &deviceRepo, &deviceDataRepo)
	deviceDataService := service.NewDeviceDataService(&deviceRepo, &deviceDataRepo)

	userHandler := handler.NewUserHandler(&userService)
	deviceHandler := handler.NewDeviceHandler(&deviceService)
	deviceDataHandler := handler.NewDeviceDataHandler(&deviceDataService)

	return userHandler, deviceHandler, deviceDataHandler
}
