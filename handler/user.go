package handler

import (
	"net/http"
	"temperature-backend/handler/util"
	"temperature-backend/service"
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
	req := registerUserRequest{}
	err := util.ParseJsonBody(w, r, &req)
	if err != nil {
		return
	}
	newUser, serviceError := (*h.service).Register(req.Email)
	if serviceError != nil {
		util.HandleServiceError(w, *serviceError)
		return
	}
	resp := registerUserResponse{Token: newUser.Token}
	util.SetResponse(w, resp, http.StatusCreated)
}

func NewUserHandler(userService *service.UserService) UserHandler {
	handler := UserHandler{service: userService}
	return handler
}
