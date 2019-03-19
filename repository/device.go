package repository

import (
	"errors"
	"sync"
	"temperature-backend/model"
)

type DeviceRepository interface {
	Save(device *model.Device) error
	FindByUserId(userId int) []model.Device
	FindByToken(token string) *model.Device
	FindById(id int) *model.Device
}

type inMemoryDeviceRepository struct {
	mux      sync.Mutex
	devices  *[]*model.Device
	sequence *seq
}

func (repo inMemoryDeviceRepository) findById(id int) *model.Device {
	for _, device := range *repo.devices {
		if device.Id == id {
			return device
		}
	}
	return nil
}

func (repo inMemoryDeviceRepository) Save(device *model.Device) error {
	repo.mux.Lock()
	defer repo.mux.Unlock()
	if device.Id == 0 {
		repo.sequence.val++
		device.Id = repo.sequence.val
		newDevice := *device
		*repo.devices = append(*repo.devices, &newDevice)
	} else {
		deviceToUpdate := repo.findById(device.Id)
		if deviceToUpdate == nil {
			return errors.New("can't find device to update")
		}
		*deviceToUpdate = *device
	}
	return nil
}

func (repo inMemoryDeviceRepository) FindByUserId(userId int) []model.Device {
	var result []model.Device
	for _, device := range *repo.devices {
		if device.UserId == userId {
			result = append(result, *device)
		}
	}
	return result
}

func (repo inMemoryDeviceRepository) FindByToken(token string) *model.Device {
	for _, device := range *repo.devices {
		if device.Token == token {
			return device
		}
	}
	return nil
}

func (repo inMemoryDeviceRepository) FindById(id int) *model.Device {
	for _, device := range *repo.devices {
		if device.Id == id {
			return device
		}
	}
	return nil
}

func NewDeviceRepository() DeviceRepository {
	devices := new([]*model.Device)
	return inMemoryDeviceRepository{devices: devices, sequence: new(seq)}
}
