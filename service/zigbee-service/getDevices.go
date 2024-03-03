package zigbee_service

// grpcurl -plaintext localhost:9080 zigbee_service.ZigbeeService/GetDevices

import (
	"context"
	"fmt"
	"log"

	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/domain"
	zigbee_service_pb "github.com/vadim-dmitriev/mirror-zigbee-api/pkg/zigbee-service"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (zs *zigbeeService) GetDevices(ctx context.Context, request *zigbee_service_pb.Empty) (*zigbee_service_pb.GetDevicesResponse, error) {
	log.Printf("'GetDevices' executed\n")

	devices, err := zs.zigbeeClient.GetAllDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all devices: %w", err)
	}

	pbDevices := make([]*zigbee_service_pb.Device, 0, len(devices))
	for _, device := range devices {
		pbDevice := mapDevicePb(device)

		pbDevices = append(pbDevices, pbDevice)
	}

	return &zigbee_service_pb.GetDevicesResponse{
		Devices: pbDevices,
	}, nil
}

func mapDevicePb(device *domain.Device) *zigbee_service_pb.Device {
	fmt.Println(device.Name)
	readables := make([]*zigbee_service_pb.State, 0)
	for _, readable := range device.Readable {
		readables = append(readables, mapStatePb(readable))
	}

	editables := make([]*zigbee_service_pb.State, 0)
	for _, editable := range device.Editable {
		editables = append(editables, mapStatePb(editable))
	}

	pbDevice := &zigbee_service_pb.Device{
		Name: device.Name,
		Characteristics: &zigbee_service_pb.Device_Characteristics{
			Description: device.Characteristics.Description,
			Vendor:      device.Characteristics.Vendor,
			Model:       device.Characteristics.Model,
		},
		Reachable: wrapperspb.Bool(device.Reachable),
		Readable:  readables,
		Editable:  editables,
	}

	return pbDevice
}
