package service

import (
	"fmt"
	"temperature-backend/model"
	"temperature-backend/repository"
	"temperature-backend/service/err"
)

type DeviceService interface {
	Register(deviceName string, userEmail string) (*model.Device, *Error)
}

type deviceServiceImpl struct {
	deviceRepo     *repository.DeviceRepository
	userRepo       *repository.UserRepository
	tokenGenerator tokenGenerator
}

func (s deviceServiceImpl) Register(deviceName string, userEmail string) (*model.Device, *Error) {
	user := (*s.userRepo).FindByEmail(userEmail)
	if user == nil {
		msg := fmt.Sprintf("user with email %s not found", userEmail)
		return nil, &Error{Code: err.UserNotFound, Msg: msg}
	}
	token := ""
	for {
		generatedToken, e := s.tokenGenerator.getToken()
		if e != nil {
			return nil, &Error{Code: UnexpectedError, Msg: e.Error()}
		}
		if (*s.deviceRepo).FindByToken(generatedToken) == nil {
			token = generatedToken
			break
		}
	}
	device := model.Device{UserId: user.Id, Name: deviceName, Token: token}
	e := (*s.deviceRepo).Save(&device)
	if e != nil {
		return nil, &Error{Code: UnexpectedError, Msg: e.Error()}
	}
	return &device, nil
}
