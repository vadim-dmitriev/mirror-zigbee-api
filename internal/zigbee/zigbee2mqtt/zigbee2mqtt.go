package zigbee2mqtt

// zigbee2mqtt topics: https://www.zigbee2mqtt.io/guide/usage/mqtt_topics_and_messages.html

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/domain"
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/zigbee"
)

const (
	zigbee2mqttDSN = "mqtt://localhost:1883/"
	clientName     = "mirror-zigbee-api"

	setDeviceStateTopicTemplate = "zigbee2mqtt/%s/set"
)

type zigbee2mqtt struct {
	mqttClient mqtt.Client

	devices []*device
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
	var devices []*device

	if err := json.Unmarshal(message.Payload(), &devices); err != nil {
		panic(err)
	}

	z2m.devices = devices
}

func (z2m *zigbee2mqtt) GetAllDevices(ctx context.Context) ([]*domain.Device, error) {
	devices := make([]*domain.Device, 0, len(z2m.devices))

	for _, rawDevice := range z2m.devices {
		device := &domain.Device{
			Name: rawDevice.FriendlyName,
			Characteristics: &domain.Characteristics{
				Description: rawDevice.Definition.Description,
				Vendor:      rawDevice.Definition.Vendor,
				Model:       rawDevice.Definition.Model,
			},
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (z2m *zigbee2mqtt) SetDeviceStatus(ctx context.Context, name domain.DeviceName, status domain.DeviceStatus) error {
	topic := fmt.Sprintf(setDeviceStateTopicTemplate, name)

	message := map[string]string{
		"state": deviceStatus[status],
	}

	marshalledMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshall message: %w", err)
	}

	log.Printf("publish message '%s' into '%s'", marshalledMessage, topic)

	token := z2m.mqttClient.Publish(topic, byte(2), true, marshalledMessage)

	<-token.Done()

	return nil
}

// func (z2m *zigbee2mqtt) Serve(devices ...zigbee.Device) error {
// 	tokens := make([]mqtt.Token, 0, len(devices))

// 	for _, device := range devices {
// 		// ИНтерфейс Device наследует zigbee.Device - поэтому кастим.
// 		castedDevice := device.(Device)

// 		log.Printf("device with topic '%s' added\n", castedDevice.GetTopic())

// 		token := z2m.mqttClient.Subscribe(castedDevice.GetTopic(), byte(0), castedDevice.GetHandler())

// 		tokens = append(tokens, token)
// 	}

// 	for _, tocken := range tokens {
// 		go tocken.Wait()
// 	}

// 	log.Printf("zigbee2mqtt started...\n")

// 	return nil
// }
