package service

import (
	"fmt"
	"temperature-backend/model"
	"temperature-backend/repository"
	"temperature-backend/view"
	"time"
)

type DeviceService interface {
	Register(deviceName string, userEmail string) (*model.Device, *Error)
	GetList(token string) (view.DeviceListView, *Error)
	GetDataList(token string, id int) (view.DeviceDataView, *Error)
}

type deviceServiceImpl struct {
	deviceRepo     *repository.DeviceRepository
	userRepo       *repository.UserRepository
	deviceDataRepo *repository.DeviceDataRepository
	tokenGenerator tokenGenerator
}

func (s deviceServiceImpl) Register(deviceName string, userEmail string) (*model.Device, *Error) {
	user := (*s.userRepo).FindByEmail(userEmail)
	if user == nil {
		msg := fmt.Sprintf("user with email %s not found", userEmail)
		return nil, &Error{Code: EntityNotFound, Msg: msg}
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

func (s deviceServiceImpl) GetList(token string) (view.DeviceListView, *Error) {
	user := (*s.userRepo).FindByToken(token)
	if user == nil {
		msg := fmt.Sprintf("can't authorize user with token %s", token)
		return view.DeviceListView{}, &Error{Code: AuthError, Msg: msg}
	}
	devices := (*s.deviceRepo).FindByUserId(user.Id)
	result := view.DeviceListView{Devices: []view.DeviceListItem{}}
	for _, device := range devices {
		data := (*s.deviceDataRepo).FindByDeviceID(device.Id)
		dateTime := ""
		if len(data) != 0 {
			dateTime = time.Unix(data[len(data)-1].Timestamp, 0).Format(time.RFC3339)
		}
		result.Devices = append(result.Devices, view.DeviceListItem{DeviceId: device.Id, DeviceName: device.Name, LastDataTime: dateTime})
	}
	return result, nil
}

func (s deviceServiceImpl) GetDataList(token string, id int) (view.DeviceDataView, *Error) {
	user := (*s.userRepo).FindByToken(token)
	if user == nil {
		msg := fmt.Sprintf("can't authorize user with token %s", token)
		return view.DeviceDataView{}, &Error{Code: AuthError, Msg: msg}
	}
	device := (*s.deviceRepo).FindById(id)
	if device == nil {
		msg := fmt.Sprintf("can't find device with id %d", id)
		return view.DeviceDataView{}, &Error{Code: AuthError, Msg: msg}
	}
	result := view.DeviceDataView{}
	deviceData := (*s.deviceDataRepo).FindByDeviceID(device.Id)
	for _, data := range deviceData {
		dateTime := time.Unix(data.Timestamp, 0).Format(time.RFC3339)
		result.Data = append(result.Data, view.DeviceDataItem{DateTime: dateTime, Temperature: data.Temperature})
	}
	return result, nil
}
func NewDeviceService(userRepo *repository.UserRepository, deviceRepo *repository.DeviceRepository, deviceDataRepo *repository.DeviceDataRepository) DeviceService {
	return deviceServiceImpl{deviceRepo: deviceRepo, userRepo: userRepo, deviceDataRepo: deviceDataRepo, tokenGenerator: tokenGeneratorImpl{}}
}
