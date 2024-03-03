package zigbee

import (
	"context"

	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/domain"
)

// Zigbee интерфейс взаимодействия с устройствами умного дома
// работающий по Zigbee протоколу.
type Zigbee interface {
	GetAllDevices(context.Context) ([]*domain.Device, error)
	SetDeviceState(context.Context, string, []*domain.DeviceState) ([]*domain.DeviceState, error)
}
