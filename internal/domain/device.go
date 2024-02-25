package domain

type DeviceName = string

type DeviceStatus = bool

type Device struct {
	Name            DeviceName
	Characteristics *Characteristics
}

type Characteristics struct {
	Description string `json:"description"`
	Vendor      string `json:"vendor"`
	Model       string `json:"model"`
}
