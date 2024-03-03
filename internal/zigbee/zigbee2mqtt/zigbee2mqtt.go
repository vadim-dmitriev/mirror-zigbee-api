package zigbee2mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/domain"
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/zigbee"
)

const (
	zigbee2mqttDSN = "mqtt://localhost:1883/"
	clientName     = "mirror-zigbee-api"

	// zigbee2mqtt topics: https://www.zigbee2mqtt.io/guide/usage/mqtt_topics_and_messages.html
	getDevicesDefinitionsTopic      = "zigbee2mqtt/bridge/devices"
	deviceStatesTopicTemplate       = "zigbee2mqtt/%s"
	deviceAvailabilityTopicTemplate = "zigbee2mqtt/%s/availability"
	setDeviceStateTopicTemplate     = "zigbee2mqtt/%s/set"
	getDeviceStateTopicTemplate     = "zigbee2mqtt/%s/get"
)

type zigbee2mqtt struct {
	mqttClient mqtt.Client

	devicesDefinitions  []*device
	devicesAvailability map[string]bool
	devicesStates       map[string]deviceStates
}

// New создает объект, имплементирующий интерфейс zigbee.Zigbee
// основанный на проекте zigbee2mqtt.
func New() (zigbee.Zigbee, error) {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(zigbee2mqttDSN)
	opts.SetClientID(clientName)

	mqttClient := mqtt.NewClient(opts)

	// Connecting to mqtt broker.
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to mqtt: %w", token.Error())
	}
	log.Printf("[DONE] connecting to mqtt.")

	// Getting all devices definitions.
	devicesDefinitions, err := getDevicesDefinitions(mqttClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices definitions: %w", err)
	}
	log.Printf("[DONE] getting devices definitions. Found %d devices.\n", len(devicesDefinitions))

	z2m := &zigbee2mqtt{
		mqttClient:          mqttClient,
		devicesDefinitions:  devicesDefinitions,
		devicesAvailability: make(map[string]bool, len(devicesDefinitions)),
		devicesStates:       make(map[string]map[string]interface{}, len(devicesDefinitions)),
	}

	devicesAvailabilityTopics := make(map[string]byte, len(devicesDefinitions))
	for _, device := range devicesDefinitions {
		deviceAvailabilityTopic := fmt.Sprintf(deviceAvailabilityTopicTemplate, device.FriendlyName)
		devicesAvailabilityTopics[deviceAvailabilityTopic] = byte(0)
	}
	if token := mqttClient.SubscribeMultiple(devicesAvailabilityTopics, z2m.deviceAvailabilityHandler); token.Error() != nil {
		return nil, fmt.Errorf("failed to multiple subscribe to device availability topics: %w", token.Error())
	}

	deviceStatesTopics := make(map[string]byte, len(devicesDefinitions))
	for _, device := range devicesDefinitions {
		deviceStatesTopic := fmt.Sprintf(deviceStatesTopicTemplate, device.FriendlyName)
		deviceStatesTopics[deviceStatesTopic] = byte(0)
	}
	if token := mqttClient.SubscribeMultiple(deviceStatesTopics, z2m.deviceStatesHandler); token.Error() != nil {
		return nil, fmt.Errorf("failed to multiple subscribe to device states topics: %w", token.Error())
	}

	return z2m, nil
}

func (z2m *zigbee2mqtt) deviceAvailabilityHandler(client mqtt.Client, message mqtt.Message) {
	parts := strings.Split(message.Topic(), "/")
	deviceName := parts[1]

	payload := make(map[string]string, 1)
	if err := json.Unmarshal(message.Payload(), &payload); err != nil {
		log.Printf("failed to unmarshal availaility payload '%s' for device '%s': %s\n", string(message.Payload()), deviceName, err)
		return
	}

	availability := false
	if payload["state"] == "online" {
		availability = true
	}

	z2m.devicesAvailability[deviceName] = availability

	log.Printf("device '%s' availability is '%s'\n",
		deviceName,
		map[bool]string{
			true:  "ONLINE",
			false: "OFFLINE",
		}[availability],
	)
}

func (z2m *zigbee2mqtt) deviceStatesHandler(client mqtt.Client, message mqtt.Message) {
	parts := strings.Split(message.Topic(), "/")
	deviceName := parts[1]

	payload := make(map[string]interface{})
	if err := json.Unmarshal(message.Payload(), &payload); err != nil {
		log.Printf("failed to unmarshal states payload '%s' for device '%s': %s\n", string(message.Payload()), deviceName, err)
		return
	}

	z2m.devicesStates[deviceName] = payload

	log.Printf("device '%s' has a new state '%s'\n", deviceName, string(message.Payload()))
}

func getDevicesDefinitions(mqttClient mqtt.Client) ([]*device, error) {
	getDevicesDefinitionsDoneChan := make(chan struct{})
	getDevicesDefinitionsErrorChan := make(chan error)

	defer close(getDevicesDefinitionsDoneChan)
	defer close(getDevicesDefinitionsErrorChan)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	devices := make([]*device, 0)

	mqttClient.Subscribe(getDevicesDefinitionsTopic, byte(0), func(client mqtt.Client, message mqtt.Message) {
		if err := json.Unmarshal(message.Payload(), &devices); err != nil {
			getDevicesDefinitionsErrorChan <- err
			return
		}
		getDevicesDefinitionsDoneChan <- struct{}{}
	})

	select {
	case <-getDevicesDefinitionsDoneChan:
		return devices, nil
	case err := <-getDevicesDefinitionsErrorChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (z2m *zigbee2mqtt) GetAllDevices(ctx context.Context) ([]*domain.Device, error) {
	devices := make([]*domain.Device, 0, len(z2m.devicesDefinitions))

	for _, rawDevice := range z2m.devicesDefinitions {

		readable := make([]*domain.DeviceState, 0)
		editable := make([]*domain.DeviceState, 0)
		for _, rawDeviceExpose := range rawDevice.Definition.Exposes {

			for _, rawFeature := range rawDeviceExpose.Features {
				if (rawFeature.Access>>0)&1 == 1 {
					readable = append(readable, mapExpose(rawFeature))
				}
				if (rawFeature.Access>>1)&1 == 1 {
					editable = append(editable, mapExpose(rawFeature))
				}
			}

			if (rawDeviceExpose.Access>>0)&1 == 1 {
				readable = append(readable, mapExpose(rawDeviceExpose))
			}
			if (rawDeviceExpose.Access>>1)&1 == 1 {
				editable = append(editable, mapExpose(rawDeviceExpose))
			}

		}

		device := &domain.Device{
			Name:   rawDevice.FriendlyName,
			Enable: !rawDevice.Disabled,
			Characteristics: &domain.Characteristics{
				Description: rawDevice.Definition.Description,
				Vendor:      rawDevice.Definition.Vendor,
				Model:       rawDevice.Definition.Model,
			},
			Readable: readable,
			Editable: editable,
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (z2m *zigbee2mqtt) SetDeviceState(ctx context.Context, name string, states []*domain.DeviceState) ([]*domain.DeviceState, error) {
	topic := fmt.Sprintf(setDeviceStateTopicTemplate, name)

	message := map[string]interface{}{}
	for _, state := range states {
		message[state.Name] = state.Value
	}

	marshalledMessage, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall message: %w", err)
	}

	log.Printf("publish message '%s' into '%s'", marshalledMessage, topic)

	_ = z2m.mqttClient.Publish(topic, byte(2), true, marshalledMessage)

	// <-token.Done()

	return states, nil
}
