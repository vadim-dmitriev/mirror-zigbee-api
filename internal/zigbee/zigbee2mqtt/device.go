package zigbee2mqtt

// Device описывает интерфейс взаимодействия с устройствами
// умного дома Zigbee, которые умеют работать по технологии zigbee2mqtt.
// type Device interface {
// 	// Наследует общий интерфейс устройств - zigbee.Device.
// 	zigbee.Device

// 	GetTopic() string
// 	GetHandler() mqtt.MessageHandler
// }

type Device struct {
	FriendlyName string `json:"friendly_name"`
}
