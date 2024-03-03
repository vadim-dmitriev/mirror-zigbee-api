package zigbee2mqtt

import (
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/domain"
)

type device struct {
	FriendlyName string     `json:"friendly_name"`
	Definition   definition `json:"definition"`
	Disabled     bool       `json:"disabled"`
	Type         string     `json:"type"`
}

func (d *device) isCoordinator() bool {
	return d.Type == "Coordinator"
}

func (d *device) triggerDeviceStates(mqttClient mqtt.Client) {
	topic := fmt.Sprintf(getDeviceStateTopicTemplate, d.FriendlyName)

	firstStateName := ""
	for _, expose := range d.Definition.Exposes {
		if (expose.Access>>2)&1 == 1 {
			firstStateName = expose.Name
			break
		}

		if expose.hasFeatures() {
			for _, feature := range expose.Features {
				if (feature.Access>>2)&1 == 1 {
					firstStateName = expose.Features[0].Name
					goto firstStateNameFound
				}
			}
		}
	}

firstStateNameFound:

	bytes, _ := json.Marshal(map[string]interface{}{
		firstStateName: "",
	})

	log.Printf("publish at topic '%s' message '%s'", topic, string(bytes))

	mqttClient.Publish(topic, byte(0), false, string(bytes))
}

type definition struct {
	Description string `json:"description"`
	Vendor      string `json:"vendor"`
	Model       string `json:"model"`

	// https://www.zigbee2mqtt.io/guide/usage/exposes.html
	Exposes []*expose `json:"exposes"`
}

type expose struct {
	Type        string    `json:"type"`
	Features    []*expose `json:"features,omitempty"`
	Access      uint      `json:"access"`
	Name        string    `json:"name"`
	Property    string    `json:"property"`
	Description string    `json:"description"`
}

func (e *expose) hasFeatures() bool {
	return len(e.Features) != 0
}

type deviceStates = map[string]interface{}

func mapExpose(expose *expose, value interface{}) *domain.DeviceState {
	exposeType := domain.DeviceStateType_UNKNOWN

	switch expose.Type {
	case "binary":
		exposeType = domain.DeviceStateType_BOOLEAN
	case "numeric":
		exposeType = domain.DeviceStateType_INTEGER
	case "text":
		exposeType = domain.DeviceStateType_STRING

	}

	e := &domain.DeviceState{
		Type:  exposeType,
		Name:  expose.Name,
		Value: value,
	}

	return e
}
