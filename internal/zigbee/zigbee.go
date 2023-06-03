package zigbee

// Zigbee интерфейс взаимодействия с устройствами умного дома
// работающий по Zigbee протоколу.
type Zigbee interface {
	Serve(devices ...Device) error
}
