package repository

import (
	"errors"
	"sync"
	"temperature-backend/model"
)

type DeviceRepository interface {
	Save(device *model.Device) error
	FindByUserId(userId int) *model.Device
	FindByToken(token string) *model.Device
}

type inMemoryDeviceRepository struct {
	mux      sync.Mutex
	devices  *[]*model.Device
	sequence int
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
		repo.sequence++
		device.Id = repo.sequence
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

func (repo inMemoryDeviceRepository) FindByUserId(userId int) *model.Device {
	for _, device := range *repo.devices {
		if device.UserId == userId {
			return device
		}
	}
	return nil
}

func (repo inMemoryDeviceRepository) FindByToken(token string) *model.Device {
	for _, device := range *repo.devices {
		if device.Token == token {
			return device
		}
	}
	return nil
}

func NewDeviceRepository() DeviceRepository {
	devices := new([]*model.Device)
	return inMemoryDeviceRepository{devices: devices}
}
