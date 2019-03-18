package repository

import (
	"errors"
	"sort"
	"sync"
	"temperature-backend/model"
)

type DeviceDataRepository interface {
	Save(deviceData *model.DeviceData) error
	FindByDeviceID(deviceId int) []model.DeviceData
}
type inMemoryDeviceDataRepo struct {
	mux      sync.Mutex
	data     map[int][]*model.DeviceData
	sequence *seq
}

func (r inMemoryDeviceDataRepo) Save(deviceData *model.DeviceData) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if deviceData.Id != 0 {
		return errors.New("deviceData.id is already set")
	}
	r.sequence.val++
	deviceData.Id = r.sequence.val
	newData := *deviceData
	r.data[deviceData.DeviceId] = append(r.data[deviceData.DeviceId], &newData)
	sort.Slice(r.data[deviceData.DeviceId], func(i, j int) bool {
		return r.data[deviceData.DeviceId][i].Timestamp < r.data[deviceData.DeviceId][j].Timestamp
	})
	return nil
}

/*
Returns device data by deviceId sorted in ascending order of timestamp field
*/
func (r inMemoryDeviceDataRepo) FindByDeviceID(deviceId int) []model.DeviceData {
	var result []model.DeviceData
	for _, d := range r.data[deviceId] {
		if d.DeviceId == deviceId {
			result = append(result, *d)
		}
	}
	return result
}

func NewDeviceDataRepository() DeviceDataRepository {
	data := make(map[int][]*model.DeviceData, 0)
	return inMemoryDeviceDataRepo{data: data, sequence: new(seq)}
}
