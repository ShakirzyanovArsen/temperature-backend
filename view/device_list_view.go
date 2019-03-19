package view

type DeviceListItem struct {
	DeviceId     int    `json:"device_id"`
	DeviceName   string `json:"device_name"`
	LastDataTime string `json:"last_data_time"`
}

type DeviceListView struct {
	Devices []DeviceListItem `json:"devices"`
}
