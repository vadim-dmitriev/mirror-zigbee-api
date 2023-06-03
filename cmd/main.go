package main

import (
	"fmt"
	"os"
	"os/signal"

	ledapi "github.com/vadim-dmitriev/mirror-zigbee-api/internal/clients/led_api"
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/zigbee/zigbee2mqtt"
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/zigbee/zigbee2mqtt/devices"
)

const (
	zigbee2mqttDSN = "mqtt://localhost:1883/"

	ledAPIHost = "localhost:8090"
)

func main() {
	// Создаем всех клиентов.
	ledAPIClent := ledapi.New(ledAPIHost)

	// Создаем zigbee контроллер.
	z2mClient := zigbee2mqtt.New(zigbee2mqttDSN)

	// Добавляем устройства.
	button1 := devices.NewWxkg01lm(ledAPIClent)

	// Стартуем zigbee котроллер.
	if err := z2mClient.Serve(button1); err != nil {
		panic(err)
	}

	// Программа работает пока не будет завершения.
	waitForInterrupt()
}

func waitForInterrupt() {
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt)

	select {
	case <-exitChan:
		fmt.Println("exiting...")
	}
}
