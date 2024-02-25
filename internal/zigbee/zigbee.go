package zigbee

import "context"

// Zigbee интерфейс взаимодействия с устройствами умного дома
// работающий по Zigbee протоколу.
type Zigbee interface {
	Serve(devices ...Device) error

	GetAllDevices(context.Context) ([]*Device, error)
}

type Device struct {
	Name string
}
