package handler

import (
	"net/http"
	"temperature-backend/handler/util"
	"temperature-backend/service"
)

type DeviceDataHandler struct {
	service *service.DeviceDataService
}
type pushDataRequest struct {
	DateTime    string  `json:"date_time"`
	Temperature float64 `json:"temperature"`
}

type pushDataResponse struct {
	Saved bool `json:"saved"`
}

func (h DeviceDataHandler) PushData(w http.ResponseWriter, r *http.Request) {
	req := pushDataRequest{}
	err := util.ParseJsonBody(w, r, &req)
	if err != nil {
		return
	}
	token := r.Header.Get("Authorization")
	serviceError := (*h.service).Push(token, req.DateTime, req.Temperature)
	if serviceError != nil {
		util.HandleServiceError(w, *serviceError)
		return
	}
	util.SetResponse(w, pushDataResponse{Saved: true})
}

func NewDeviceDataHandler(service *service.DeviceDataService) DeviceDataHandler {
	return DeviceDataHandler{service: service}
}
