package main

import (
	"context"
	"fmt"

	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/server"
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/zigbee/zigbee2mqtt"
	zigbee_service_pb "github.com/vadim-dmitriev/mirror-zigbee-api/pkg/zigbee-service"
	zigbee_service "github.com/vadim-dmitriev/mirror-zigbee-api/service/zigbee-service"
)

const (
	// grpcPort порт для gRPC соединений.
	grpcPort = 9080
	// httpPort порт для HTTP соединений.
	httpPort = 9081
)

func main() {
	ctx := context.Background()

	zigbeeClient, err := zigbee2mqtt.New()
	if err != nil {
		panic(err)
	}

	zigbeeService, err := zigbee_service.New(
		zigbeeClient,
	)
	if err != nil {
		panic(err)
	}

	grpcServer, err := server.NewGRPC(grpcPort)
	if err != nil {
		panic(fmt.Errorf("failed to create gRPS server: %w", err))
	}
	grpcServer.RegisterService(&zigbee_service_pb.ZigbeeService_ServiceDesc, zigbeeService)

	go grpcServer.Run()

	httpServer, err := server.NewHTTP(httpPort, grpcPort)
	if err != nil {
		panic(fmt.Errorf("failed to create HTTP server: %w", err))
	}

	httpServer.Run(ctx)
}
