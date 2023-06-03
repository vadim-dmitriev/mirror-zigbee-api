package devices

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	ledapi "github.com/vadim-dmitriev/mirror-zigbee-api/internal/clients/led_api"
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/zigbee/zigbee2mqtt"
)

const (
	wxkg01lmTopic = "zigbee2mqtt/button-1"
)

type wxkg01lm struct {
	topic string

	ledAPIClient ledapi.LedAPI
}

type wxkg01lmMessage struct {
	// Battery Remaining battery in %, can take up to 24 hours before reported.
	Battery int `json:"battery"`
	// Linkquality Качество связи (мощность сигнала)
	Linkquality int `json:"linkquality"`
	// PowerOutageCount Количество отключений электроэнергии
	PowerOutageCount int `json:"power_outage_count"`
	// Voltage Напряжение аккумулятора в милливольтах
	Voltage int `json:"voltage"`
	// Click Инициированное действие (например, нажатие кнопки)
	Action string `json:"action"`
}

// NewWxkg01lm конструктор умного zigbee устройства: "Mi Wireless Switch".
func NewWxkg01lm(ledAPIClient ledapi.LedAPI) zigbee2mqtt.Device {
	return &wxkg01lm{
		topic: wxkg01lmTopic,

		ledAPIClient: ledAPIClient,
	}
}

func (w *wxkg01lm) GetTopic() string {
	return w.topic
}

func (w *wxkg01lm) GetHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, message mqtt.Message) {
		msg := wxkg01lmMessage{}

		if err := json.Unmarshal(message.Payload(), &msg); err != nil {
			panic(err)
		}

		if msg.Action == "single" {
			if err := w.ledAPIClient.SwitchLED(); err != nil {
				panic(err)
			}
		}

		log.Printf("Message recieved: %s\n", string(message.Payload()))
	}
}
