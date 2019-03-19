package service

import (
	"fmt"
	"reflect"
	"temperature-backend/model"
	"temperature-backend/repository"
	"temperature-backend/view"
	"testing"
	"time"
)

type deviceRepoMock struct {
	saveError          error
	findByUserIdResult *model.Device
	findByTokenResult  *model.Device
	findByIdResult     *model.Device
}

func (r deviceRepoMock) Save(device *model.Device) error {
	if r.saveError == nil {
		device.Id = 1
		return nil
	}
	return r.saveError
}

func (r deviceRepoMock) FindByUserId(userId int) []model.Device {
	return []model.Device{*r.findByUserIdResult}
}

func (r deviceRepoMock) FindByToken(token string) *model.Device {
	return r.findByTokenResult
}

func (r deviceRepoMock) FindById(id int) *model.Device {
	return r.findByIdResult
}

func createDeviceRepoMock(err error, findByUserIdResult *model.Device, findByTokenResult *model.Device,
	findByIdResult *model.Device) repository.DeviceRepository {
	res := deviceRepoMock{saveError: err,
		findByUserIdResult: findByUserIdResult,
		findByTokenResult:  findByTokenResult,
		findByIdResult:     findByIdResult,
	}
	return res
}

func Test_deviceServiceImpl_Register(t *testing.T) {
	type fields struct {
		deviceRepo     *repository.DeviceRepository
		userRepo       *repository.UserRepository
		tokenGenerator tokenGenerator
	}
	type args struct {
		deviceName string
		userEmail  string
	}

	findByEmailUserRepo := createUserRepoMock(nil, &model.User{Id: 1, Email: "test@test.ru"}, nil)
	userNotFoundRepo := createUserRepoMock(nil, nil, nil)

	successDeviceRepo := createDeviceRepoMock(nil, nil, nil, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Device
		wantErr *Error
	}{
		{
			name: "success register",
			fields: fields{
				deviceRepo:     &successDeviceRepo,
				userRepo:       &findByEmailUserRepo,
				tokenGenerator: creatMockedTokenGenerator([]string{"123"}),
			},
			args:    args{deviceName: "device1", userEmail: "test@test.ru"},
			want:    &model.Device{Id: 1, UserId: 1, Name: "device1", Token: "123"},
			wantErr: nil,
		},
		{
			name: "user not found",
			fields: fields{
				deviceRepo:     &successDeviceRepo,
				userRepo:       &userNotFoundRepo,
				tokenGenerator: creatMockedTokenGenerator([]string{"123"}),
			},
			args:    args{deviceName: "device1", userEmail: "test@test.ru"},
			want:    nil,
			wantErr: &Error{Code: EntityNotFound, Msg: "user with email test@test.ru not found"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := deviceServiceImpl{
				deviceRepo:     tt.fields.deviceRepo,
				userRepo:       tt.fields.userRepo,
				tokenGenerator: tt.fields.tokenGenerator,
			}
			got, got1 := s.Register(tt.args.deviceName, tt.args.userEmail)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deviceServiceImpl.Register() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.wantErr) {
				t.Errorf("deviceServiceImpl.Register() got1 = %v, want %v", got1, tt.wantErr)
			}
		})
	}
}

func Test_deviceServiceImpl_GetList(t *testing.T) {
	type fields struct {
		deviceRepo     *repository.DeviceRepository
		userRepo       *repository.UserRepository
		deviceDataRepo *repository.DeviceDataRepository
		tokenGenerator tokenGenerator
	}
	type args struct {
		token string
	}
	validToken := "validToken"
	invalidToken := "invalidToken"
	successDeviceRepoMock := createDeviceRepoMock(nil, &model.Device{Id: 2, Name: "deviceName1"}, nil, nil)
	successUserRepoMock := createUserRepoMock(nil, nil, &model.User{Id: 1})
	errorUserRepoMock := createUserRepoMock(nil, nil, nil)
	successDataRepoMock := createDeviceDataRepoMock(nil, []model.DeviceData{
		{Id: 1, Timestamp: 1552942800, Temperature: 30.0},
		{Id: 2, Timestamp: 1552946400, Temperature: 20.0},
	})
	generator := creatMockedTokenGenerator([]string{})
	tests := []struct {
		name   string
		fields fields
		args   args
		want   view.DeviceListView
		wantE  *Error
	}{
		{
			name:   "success",
			fields: fields{deviceRepo: &successDeviceRepoMock, userRepo: &successUserRepoMock, deviceDataRepo: &successDataRepoMock, tokenGenerator: generator},
			args:   args{token: validToken},
			want: view.DeviceListView{
				Devices: []view.DeviceListItem{
					{DeviceId: 2, DeviceName: "deviceName1", LastDataTime: time.Unix(1552946400, 0).Format(time.RFC3339)},
				},
			},
			wantE: nil,
		},
		{
			name:   "user auth error",
			fields: fields{deviceRepo: &successDeviceRepoMock, userRepo: &errorUserRepoMock, deviceDataRepo: &successDataRepoMock, tokenGenerator: generator},
			args:   args{token: invalidToken},
			want:   view.DeviceListView{},
			wantE:  &Error{Code: AuthError, Msg: fmt.Sprintf("can't authorize user with token %s", invalidToken)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := deviceServiceImpl{
				deviceRepo:     tt.fields.deviceRepo,
				userRepo:       tt.fields.userRepo,
				deviceDataRepo: tt.fields.deviceDataRepo,
				tokenGenerator: tt.fields.tokenGenerator,
			}
			got, got1 := s.GetList(tt.args.token)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deviceServiceImpl.GetList() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.wantE) {
				t.Errorf("deviceServiceImpl.GetList() got1 = %v, want %v", got1, tt.wantE)
			}
		})
	}
}

func Test_deviceServiceImpl_GetDataList(t *testing.T) {
	type fields struct {
		deviceRepo     *repository.DeviceRepository
		userRepo       *repository.UserRepository
		deviceDataRepo *repository.DeviceDataRepository
		tokenGenerator tokenGenerator
	}
	type args struct {
		token string
		id    int
	}
	validToken := "validToken"
	invalidToken := "invalidToken"
	successDeviceRepoMock := createDeviceRepoMock(nil, nil, nil, &model.Device{Id: 2, Name: "deviceName1"})
	cantFindDeviceRepoMock := createDeviceRepoMock(nil, nil, nil, nil)
	successUserRepoMock := createUserRepoMock(nil, nil, &model.User{Id: 1})
	errorUserRepoMock := createUserRepoMock(nil, nil, nil)
	successDataRepoMock := createDeviceDataRepoMock(nil, []model.DeviceData{
		{Id: 1, Timestamp: 1552942800, Temperature: 30.0},
		{Id: 1, Timestamp: 1552946400, Temperature: 20.0},
	})
	tests := []struct {
		name   string
		fields fields
		args   args
		want   view.DeviceDataView
		wantE  *Error
	}{
		{
			name:   "success",
			fields: fields{deviceRepo: &successDeviceRepoMock, userRepo: &successUserRepoMock, deviceDataRepo: &successDataRepoMock},
			args:   args{token: validToken, id: 2},
			want: view.DeviceDataView{Data: []view.DeviceDataItem{
				{DateTime: "2019-03-19T02:00:00+05:00", Temperature: 30.},
				{DateTime: "2019-03-19T03:00:00+05:00", Temperature: 20.},
			}},
			wantE: nil,
		},
		{
			name:   "cant authorize user token",
			fields: fields{deviceRepo: &successDeviceRepoMock, userRepo: &errorUserRepoMock, deviceDataRepo: &successDataRepoMock},
			args:   args{token: invalidToken, id: 2},
			want:   view.DeviceDataView{},
			wantE:  &Error{Code: AuthError, Msg: fmt.Sprintf("can't authorize user with token %s", invalidToken)},
		},
		{
			name:   "device not found",
			fields: fields{deviceRepo: &cantFindDeviceRepoMock, userRepo: &successUserRepoMock, deviceDataRepo: &successDataRepoMock},
			args:   args{token: validToken, id: 100},
			want:   view.DeviceDataView{},
			wantE:  &Error{Code: AuthError, Msg: fmt.Sprintf("can't find device with id %d", 100)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := deviceServiceImpl{
				deviceRepo:     tt.fields.deviceRepo,
				userRepo:       tt.fields.userRepo,
				deviceDataRepo: tt.fields.deviceDataRepo,
				tokenGenerator: tt.fields.tokenGenerator,
			}
			got, got1 := s.GetDataList(tt.args.token, tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deviceServiceImpl.GetList() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.wantE) {
				t.Errorf("deviceServiceImpl.GetList() got1 = %v, want %v", got1, tt.wantE)
			}
		})
	}
}
