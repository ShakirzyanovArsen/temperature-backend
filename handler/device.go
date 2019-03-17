package handler

import (
	"net/http"
	"temperature-backend/handler/util"
	"temperature-backend/service"
)

type DeviceHandler struct {
	service *service.DeviceService
}

type registerDeviceRequest struct {
	DeviceName string `json:"device_name"`
	UserEmail  string `json:"user_email"`
}

type registerDeviceResponse struct {
	Token string `json:"token"`
}

func (h DeviceHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	req := registerDeviceRequest{}
	err := util.ParseJsonBody(w, r, &req)
	if err != nil {
		return
	}
	newDevice, serviceError := (*h.service).Register(req.DeviceName, req.UserEmail)
	if serviceError != nil {
		util.HandleServiceError(w, *serviceError)
		return
	}
	resp := registerDeviceResponse{Token: newDevice.Token}
	util.SetResponse(w, resp)
}

func NewDeviceHandler(service *service.DeviceService) DeviceHandler {
	return DeviceHandler{service: service}
}
