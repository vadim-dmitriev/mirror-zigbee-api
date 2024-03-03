package domain

type DeviceStateType uint8

const (
	DeviceStateType_UNKNOWN = DeviceStateType(iota)
	DeviceStateType_BOOLEAN
	DeviceStateType_INTEGER
	DeviceStateType_STRING
)

type DeviceState struct {
	Type  DeviceStateType
	Name  string
	Value interface{}
}

type Device struct {
	Name            string
	Enable          bool
	Characteristics *Characteristics
	Readable        []*DeviceState
	Editable        []*DeviceState
}

type Characteristics struct {
	Description string `json:"description"`
	Vendor      string `json:"vendor"`
	Model       string `json:"model"`
}
