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
	util.SetResponse(w, resp, http.StatusCreated)
}

func (h DeviceHandler) GetDeviceList(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	deviceListView, e := (*h.service).GetList(token)
	if e != nil {
		util.HandleServiceError(w, *e)
		return
	}
	util.SetResponse(w, deviceListView, http.StatusOK)
}

func NewDeviceHandler(service *service.DeviceService) DeviceHandler {
	return DeviceHandler{service: service}
}
