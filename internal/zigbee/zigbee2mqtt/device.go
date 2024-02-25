package zigbee2mqtt

const (
	deviceStatusOFF = "OFF"
	deviceStatusON  = "ON"
)

type device struct {
	FriendlyName string     `json:"friendly_name"`
	Definition   definition `json:"definition"`
}

type definition struct {
	Description string `json:"description"`
	Vendor      string `json:"vendor"`
	Model       string `json:"model"`
}

var deviceStatus = map[bool]string{
	false: deviceStatusOFF,
	true:  deviceStatusON,
}
