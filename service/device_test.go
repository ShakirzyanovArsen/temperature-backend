package service

import (
	"reflect"
	"temperature-backend/model"
	"temperature-backend/repository"
	"temperature-backend/service/err"
	"testing"
)

type deviceRepoMock struct {
	saveError          error
	findByUserIdResult *model.Device
	findByTokenResult  *model.Device
}

func (r deviceRepoMock) Save(device *model.Device) error {
	if r.saveError == nil {
		device.Id = 1
		return nil
	}
	return r.saveError
}

func (r deviceRepoMock) FindByUserId(userId int) *model.Device {
	return r.findByUserIdResult
}

func (r deviceRepoMock) FindByToken(token string) *model.Device {
	return r.findByTokenResult
}

func createDeviceRepoMock(err error, findByUserIdResult *model.Device, findByTokenResult *model.Device) repository.DeviceRepository {
	res := deviceRepoMock{saveError: err, findByUserIdResult: findByUserIdResult, findByTokenResult: findByTokenResult}
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

	successDeviceRepo := createDeviceRepoMock(nil, nil, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Device
		wantErr *Error
	}{
		{
			name: "succes register",
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
			wantErr: &Error{Code: err.UserNotFound, Msg: "user with email test@test.ru not found"},
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
