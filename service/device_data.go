package service

import (
	"fmt"
	"temperature-backend/model"
	"temperature-backend/repository"
	"time"
)

type DeviceDataService interface {
	Push(token string, timestamp string, temperature float64) *Error
}

type deviceDataServiceImpl struct {
	deviceRepo     *repository.DeviceRepository
	deviceDataRepo *repository.DeviceDataRepository
}

func (s deviceDataServiceImpl) Push(token string, dataTime string, temperature float64) *Error {
	device := (*s.deviceRepo).FindByToken(token)
	if device == nil {
		return &Error{Code: AuthError, Msg: fmt.Sprintf("Device auth with token %s failed", token)}
	}
	timestamp, e := time.Parse(time.RFC3339, dataTime)
	if e != nil {
		return &Error{Code: ParseError, Msg: "Can't parse date_time field. Expected date time in RFC3339 format"}
	}
	deviceData := model.DeviceData{DeviceId: device.Id, Timestamp: timestamp.Unix(), Temperature: temperature}
	e = (*s.deviceDataRepo).Save(&deviceData)
	if e != nil {
		return &Error{Code: UnexpectedError, Msg: fmt.Sprintf("error occurred while save device data: %s", e.Error())}
	}
	return nil
}

func NewDeviceDataService(deviceRepo *repository.DeviceRepository, deviceDataRepo *repository.DeviceDataRepository) DeviceDataService {
	return deviceDataServiceImpl{deviceRepo: deviceRepo, deviceDataRepo: deviceDataRepo}
}
