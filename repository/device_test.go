package repository

import (
	"temperature-backend/model"
	"temperature-backend/test_utils"
	"testing"
)

func Test_inMemoryDeviceRepo_save(t *testing.T) {
	repo := NewDeviceRepository().(inMemoryDeviceRepository)
	device := model.Device{UserId: 1, Name: "device1", Token: "token"}
	err := repo.Save(&device)

	test_utils.UnexpectedError(t, err)
	test_utils.AssertInt(t, 1, device.Id)
	test_utils.AssertInt(t, 1, len(*repo.devices))
}
func Test_inMemoryDeviceRepo_save_update(t *testing.T) {
	repo := NewDeviceRepository().(inMemoryDeviceRepository)
	device := model.Device{UserId: 1, Name: "device1", Token: "token"}
	err := repo.Save(&device)
	test_utils.UnexpectedError(t, err)
	updatedUserId := 2
	device.UserId = updatedUserId
	err = repo.Save(&device)
	test_utils.UnexpectedError(t, err)
	test_utils.AssertInt(t, updatedUserId, (*(repo.devices))[0].UserId)
}

func Test_inMemoryDeviceRepository_update_not_exists(t *testing.T) {
	repo := NewDeviceRepository().(inMemoryDeviceRepository)
	device := model.Device{Id: 1, UserId: 1, Name: "device1", Token: "token"}
	err := repo.Save(&device)
	test_utils.AssertError(t, err)
}
func Test_inMemoryDeviceRepo_findByUserId(t *testing.T) {
	repo := NewDeviceRepository().(inMemoryDeviceRepository)
	device1 := model.Device{UserId: 1, Name: "device1", Token: "token1"}
	device2 := model.Device{UserId: 1, Name: "device2", Token: "token2"}
	err := repo.Save(&device1)
	err = repo.Save(&device2)
	test_utils.UnexpectedError(t, err)
	findedDevice := repo.FindByUserId(1)
	test_utils.AssertStruct(t, &model.Device{UserId: 1, Name: "device1", Token: "token"}, findedDevice)
	findedDevice = repo.FindByUserId(2)
	test_utils.AssertInt(t, 0, len(findedDevice))
}

func Test_inMemoryDeviceRepo_findById(t *testing.T) {
	repo := NewDeviceRepository().(inMemoryDeviceRepository)
	device1 := model.Device{UserId: 1, Name: "device1", Token: "token1"}
	device2 := model.Device{UserId: 2, Name: "device2", Token: "token2"}
	err := repo.Save(&device1)
	err = repo.Save(&device2)
	test_utils.UnexpectedError(t, err)
	findedDevice := repo.FindById(1)
	test_utils.AssertStruct(t, &model.Device{Id: 1, UserId: 1, Name: "device1", Token: "token"}, findedDevice)
	findedDevice = repo.FindById(100500)
	test_utils.AssertStruct(t, nil, findedDevice)
}

func Test_inMemoryDeviceRepo_FindByToken(t *testing.T) {
	repo := NewDeviceRepository().(inMemoryDeviceRepository)
	device := model.Device{UserId: 1, Name: "device1", Token: "token"}
	err := repo.Save(&device)
	test_utils.UnexpectedError(t, err)
	findedUser := repo.FindByToken("token")
	test_utils.AssertStruct(t, []model.Device{{UserId: 1, Name: "device1", Token: "token1"}, {UserId: 1, Name: "device2", Token: "token2"}}, findedUser)
	findedUser = repo.FindByToken("123")
	test_utils.AssertStruct(t, nil, findedUser)
}

func TestNewRepository(t *testing.T) {
	repo := NewDeviceRepository().(inMemoryDeviceRepository)
	test_utils.AssertInt(t, 0, repo.sequence.val)
	test_utils.AssertNotNil(t, repo.devices)
}
