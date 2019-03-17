package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"temperature-backend/handler/handler_utils"
	"temperature-backend/service"
	"temperature-backend/util"
)

type registerUserRequest struct {
	Email string `json:"email"`
}

type registerUserResponse struct {
	Token string `json:"token"`
}

type UserHandler struct {
	service *service.UserService
}

func (h UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	reqBuf := new(bytes.Buffer)
	_, err := reqBuf.ReadFrom(r.Body)
	if err != nil {
		util.HandleInternalError(w, err)
		return
	}
	req := registerUserRequest{}
	err = json.Unmarshal(reqBuf.Bytes(), &req)
	if err != nil {
		util.HandleSerializingError(w, err)
		return
	}
	newUser, serviceError := (*h.service).Register(req.Email)
	if serviceError != nil {
		util.HandleServiceError(w, *serviceError)
		return
	}
	resp := registerUserResponse{Token: newUser.Token}
	handler_utils.SetResponse(w, resp)
}

func NewUserHandler(userService *service.UserService) UserHandler {
	handler := UserHandler{service: userService}
	return handler
}
