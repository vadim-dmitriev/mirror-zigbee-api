package zigbee2mqtt

import (
	"context"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/zigbee"
)

const (
	zigbee2mqttDSN = "mqtt://localhost:1883/"
	clientName     = "mirror-zigbee-api"
)

type zigbee2mqtt struct {
	mqttClient mqtt.Client

	devices []*Device
}

// New создает объект, имплементирующий интерфейс zigbee.Zigbee
// основанный на проекте zigbee2mqtt.
func New() (zigbee.Zigbee, error) {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(zigbee2mqttDSN)
	opts.SetClientID(clientName)

	mqttClient := mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	z2m := &zigbee2mqtt{
		mqttClient: mqttClient,
		devices:    nil,
	}

	mqttClient.Subscribe("zigbee2mqtt/bridge/devices", byte(0), z2m.devicesHandler)

	return z2m, nil
}

func (z2m *zigbee2mqtt) devicesHandler(client mqtt.Client, message mqtt.Message) {
	var devices []*Device
	if err := json.Unmarshal(message.Payload(), &devices); err != nil {
		panic(err)
	}

	z2m.devices = devices
}

func (z2m *zigbee2mqtt) Serve(devices ...zigbee.Device) error {
	// tokens := make([]mqtt.Token, 0, len(devices))

	// for _, device := range devices {
	// 	// ИНтерфейс Device наследует zigbee.Device - поэтому кастим.
	// 	castedDevice := device.(Device)

	// 	log.Printf("device with topic '%s' added\n", castedDevice.GetTopic())

	// 	token := z2m.mqttClient.Subscribe(castedDevice.GetTopic(), byte(0), castedDevice.GetHandler())

	// 	tokens = append(tokens, token)
	// }

	// for _, tocken := range tokens {
	// 	go tocken.Wait()
	// }

	// log.Printf("zigbee2mqtt started...\n")

	return nil
}

func (z2m *zigbee2mqtt) GetAllDevices(ctx context.Context) ([]*zigbee.Device, error) {
	devices := make([]*zigbee.Device, 0, len(z2m.devices))
	for _, rawDevice := range z2m.devices {
		device := &zigbee.Device{
			Name: rawDevice.FriendlyName,
		}

		devices = append(devices, device)
	}

	return devices, nil
}
