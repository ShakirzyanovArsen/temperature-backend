package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"temperature-backend/handler/handler_utils"
	"temperature-backend/service"
	"temperature-backend/util"
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
	reqBuf := new(bytes.Buffer)
	_, err := reqBuf.ReadFrom(r.Body)
	if err != nil {
		util.HandleInternalError(w, err)
		return
	}
	req := registerDeviceRequest{}
	err = json.Unmarshal(reqBuf.Bytes(), &req)
	strBody := string(reqBuf.Bytes())
	fmt.Println(strBody)
	if err != nil {
		util.HandleSerializingError(w, err)
		return
	}
	newDevice, serviceError := (*h.service).Register(req.DeviceName, req.UserEmail)
	if serviceError != nil {
		util.HandleServiceError(w, *serviceError)
		return
	}
	resp := registerDeviceResponse{Token: newDevice.Token}
	handler_utils.SetResponse(w, resp)
}

func NewDeviceHandler(service *service.DeviceService) DeviceHandler {
	return DeviceHandler{service: service}
}
