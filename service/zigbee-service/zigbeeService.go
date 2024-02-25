package zigbee_service

// grpcurl -plaintext localhost:9080 list zigbee_service.ZigbeeService

import (
	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/zigbee"
	zigbee_service_pb "github.com/vadim-dmitriev/mirror-zigbee-api/pkg/zigbee-service"
)

type zigbeeService struct {
	zigbee_service_pb.UnimplementedZigbeeServiceServer

	zigbeeClient zigbee.Zigbee
}

var _ zigbee_service_pb.ZigbeeServiceServer = &zigbeeService{}

// New создает имплементацию gRPC сервиса ZigbeeService.
func New(
	zigbeeClient zigbee.Zigbee,
) (zigbee_service_pb.ZigbeeServiceServer, error) {

	zs := &zigbeeService{
		zigbeeClient: zigbeeClient,
	}

	return zs, nil
}
