package view

type DeviceDataItem struct {
	DateTime    string  `json:"date_time"`
	Temperature float64 `json:"temperature"`
}

type DeviceDataView struct {
	Data []DeviceDataItem `json:"data"`
}
