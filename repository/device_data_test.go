package repository

import (
	"temperature-backend/model"
	"temperature-backend/test_utils"
	"testing"
)

func Test_inMemoryDeviceDataRepo_save(t *testing.T) {
	repo := NewDeviceDataRepository().(inMemoryDeviceDataRepo)
	data := model.DeviceData{DeviceId: 2, Timestamp: 100500, Temperature: 33.5}
	err := repo.Save(&data)

	test_utils.UnexpectedError(t, err)
	test_utils.AssertInt(t, 1, data.Id)
	test_utils.AssertInt(t, 1, len(repo.data[data.DeviceId]))
	test_utils.AssertStruct(t, data, (repo.data[data.DeviceId])[0])
}

func Test_inMemoryDeviceDataRepo_save_error(t *testing.T) {
	repo := NewDeviceDataRepository().(inMemoryDeviceDataRepo)
	data := model.DeviceData{Id: 1, DeviceId: 2, Timestamp: 100500, Temperature: 33.5}
	err := repo.Save(&data)
	test_utils.AssertError(t, err)
}

func Test_inMemoryDeviceDataRepo_findByDeviceId(t *testing.T) {
	repo := NewDeviceDataRepository().(inMemoryDeviceDataRepo)
	data1 := model.DeviceData{DeviceId: 2, Timestamp: 100500, Temperature: 33.5}
	data2 := model.DeviceData{DeviceId: 1, Timestamp: 800800, Temperature: 33.6}
	err := repo.Save(&data1)
	err = repo.Save(&data2)

	test_utils.UnexpectedError(t, err)
	devicDataList := repo.FindByDeviceID(2)
	test_utils.AssertInt(t, 1, len(devicDataList))
	test_utils.AssertStruct(t, []model.DeviceData{{DeviceId: 1, Timestamp: 800800, Temperature: 33.6}}, devicDataList)
}

func Test_inMemoryDeviceDataRepo_findByDeviceId_ascSort(t *testing.T) {
	repo := NewDeviceDataRepository().(inMemoryDeviceDataRepo)
	data1 := model.DeviceData{DeviceId: 1, Timestamp: 100500, Temperature: 33.5}
	data2 := model.DeviceData{DeviceId: 1, Timestamp: 800800, Temperature: 33.6}
	err := repo.Save(&data1)
	err = repo.Save(&data2)

	test_utils.UnexpectedError(t, err)
	devicDataList := repo.FindByDeviceID(1)
	test_utils.AssertInt(t, 2, len(devicDataList))
	test_utils.AssertStruct(t, []model.DeviceData{{DeviceId: 1, Timestamp: 800800, Temperature: 33.6}}, devicDataList[0])
	test_utils.AssertStruct(t, []model.DeviceData{{DeviceId: 1, Timestamp: 100500, Temperature: 33.5}}, devicDataList[1])
}
