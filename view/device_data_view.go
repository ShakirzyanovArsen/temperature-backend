package view

type DeviceDataItem struct {
	DateTime    string
	Temperature float64
}

type DeviceDataView struct {
	Data []DeviceDataItem
}
