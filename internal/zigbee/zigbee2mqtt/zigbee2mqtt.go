package zigbee2mqtt

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/zigbee"
)

const (
	clientName = "mirror-zigbee-api"
)

type zigbee2mqtt struct {
	mqtt.Client
}

// New создает объект, имплементирующий интерфейс zigbee.Zigbee
// основанный на проекте zigbee2mqtt
func New(brockerDSN string) zigbee.Zigbee {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(brockerDSN)
	opts.SetClientID(clientName)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &zigbee2mqtt{
		client,
	}
}

func (z2m *zigbee2mqtt) Serve(devices ...zigbee.Device) error {
	tokens := make([]mqtt.Token, 0, len(devices))

	for _, device := range devices {
		// ИНтерфейс Device наследует zigbee.Device - поэтому кастим.
		castedDevice := device.(Device)

		log.Printf("device with topic '%s' added\n", castedDevice.GetTopic())

		token := z2m.Client.Subscribe(castedDevice.GetTopic(), byte(0), castedDevice.GetHandler())

		tokens = append(tokens, token)
	}

	for _, tocken := range tokens {
		go tocken.Wait()
	}

	log.Printf("zigbee2mqtt started...\n")

	return nil
}
