package zigbee2mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/zigbee"
)

// Device описывает интерфейс взаимодействия с устройствами
// умного дома Zigbee, которые умеют работать по технологии zigbee2mqtt.
type Device interface {
	// Наследует общий интерфейс устройств - zigbee.Device.
	zigbee.Device

	GetTopic() string
	GetHandler() mqtt.MessageHandler
}
