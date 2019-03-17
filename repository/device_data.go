package repository

import (
	"errors"
	"sync"
	"temperature-backend/model"
)

type DeviceDataRepository interface {
	Save(deviceData *model.DeviceData) error
	FindByDeviceID(deviceId int) []model.DeviceData
}
type inMemoryDeviceDataRepo struct {
	mux      sync.Mutex
	data     *[]*model.DeviceData
	sequence int
}

func (r inMemoryDeviceDataRepo) Save(deviceData *model.DeviceData) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if deviceData.Id != 0 {
		return errors.New("deviceData.id is already set")
	}
	r.sequence++
	deviceData.Id = r.sequence
	newData := *deviceData
	*r.data = append(*r.data, &newData)
	return nil
}

func (r inMemoryDeviceDataRepo) FindByDeviceID(deviceId int) []model.DeviceData {
	var result []model.DeviceData
	for _, d := range *r.data {
		if d.DeviceId == deviceId {
			result = append(result, *d)
		}
	}
	return result
}

func NewDeviceDataRepository() DeviceDataRepository {
	data := make([]*model.DeviceData, 0)
	return inMemoryDeviceDataRepo{data: &data, sequence: 0}
}
