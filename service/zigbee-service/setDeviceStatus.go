package zigbee_service

// grpcurl -plaintext localhost:9080 zigbee_service.ZigbeeService/GetDeviceStatus

import (
	"context"
	"fmt"
	"log"

	zigbee_service_pb "github.com/vadim-dmitriev/mirror-zigbee-api/pkg/zigbee-service"
)

func (zs *zigbeeService) SetDeviceStatus(ctx context.Context, request *zigbee_service_pb.SetDeviceStatusRequest) (*zigbee_service_pb.Empty, error) {
	log.Printf("'SetDeviceStatus' executed with request '%s'\n", request)

	if err := zs.zigbeeClient.SetDeviceStatus(ctx, request.GetName(), request.GetStatus().GetValue()); err != nil {
		return nil, fmt.Errorf("failed to set device status: %w", err)
	}

	return &zigbee_service_pb.Empty{}, nil
}
