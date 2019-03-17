package service

import (
	"errors"
	"reflect"
	"temperature-backend/model"
	"temperature-backend/repository"
	"testing"
)

type deviceDataRepoMock struct {
	saveError            error
	findByDeviceIdResult []model.DeviceData
}

func (r deviceDataRepoMock) Save(deviceData *model.DeviceData) error {
	if r.saveError == nil {
		deviceData.Id = 1
	}
	return r.saveError
}

func (r deviceDataRepoMock) FindByDeviceID(deviceId int) []model.DeviceData {
	return r.findByDeviceIdResult
}

func createDeviceDataRepoMock(saveError error, findByDeviceIdResult []model.DeviceData) repository.DeviceDataRepository {
	result := &(deviceDataRepoMock{saveError: saveError, findByDeviceIdResult: findByDeviceIdResult})
	return result
}

func Test_deviceDataServiceImpl_Push(t *testing.T) {
	type fields struct {
		deviceRepo     *repository.DeviceRepository
		deviceDataRepo *repository.DeviceDataRepository
	}
	type args struct {
		token       string
		dataTime    string
		temperature float64
	}
	validToken := "123"
	successDeviceMock := createDeviceRepoMock(nil, nil, &model.Device{UserId: 1, Token: validToken})
	deviceAuthFailedMock := createDeviceRepoMock(nil, nil, nil)
	successDeviceDataMock := createDeviceDataRepoMock(nil, []model.DeviceData{})
	saveFailedDeviceDataMock := createDeviceDataRepoMock(errors.New("error"), []model.DeviceData{})
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Error
	}{
		{
			name:   "success data push",
			fields: fields{deviceRepo: &successDeviceMock, deviceDataRepo: &successDeviceDataMock},
			args:   args{token: validToken, dataTime: "2019-03-17T21:08:00+05:00", temperature: 33.6},
			want:   nil,
		},
		{
			name:   "date time parse error",
			fields: fields{deviceRepo: &successDeviceMock, deviceDataRepo: &successDeviceDataMock},
			args:   args{token: validToken, dataTime: "2019-03-17", temperature: 33.6},
			want:   &Error{Code: ParseError, Msg: "Can't parse date_time field. Expected date time in RFC3339 format"},
		},
		{
			name:   "device auth failed",
			fields: fields{deviceRepo: &deviceAuthFailedMock, deviceDataRepo: &successDeviceDataMock},
			args:   args{token: "invalidToken", dataTime: "2019-03-17T21:08:00+05:00", temperature: 33.6},
			want:   &Error{Code: AuthError, Msg: "Device auth with token invalidToken failed"},
		},
		{
			name:   "saveFailed",
			fields: fields{deviceRepo: &successDeviceMock, deviceDataRepo: &saveFailedDeviceDataMock},
			args:   args{token: validToken, dataTime: "2019-03-17T21:08:00+05:00", temperature: 33.6},
			want:   &Error{Code: UnexpectedError, Msg: "error occurred while save device data: error"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := deviceDataServiceImpl{
				deviceRepo:     tt.fields.deviceRepo,
				deviceDataRepo: tt.fields.deviceDataRepo,
			}
			if got := s.Push(tt.args.token, tt.args.dataTime, tt.args.temperature); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deviceDataServiceImpl.Push() = %v, want %v", got, tt.want)
			}
		})
	}
}
