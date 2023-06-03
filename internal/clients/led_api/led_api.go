package ledapi

import (
	"context"

	led_service "github.com/vadim-dmitriev/mirror-led-api/pkg/led-service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// LedAPI интерфейс-обертка над сервисом mirror-led-api.
type LedAPI interface {
	LightLED() error
	SwitchLED() error
}

type ledAPI struct {
	host string
}

// New создает враппер над grpc-клиентом к mirror-led-api.
func New(host string) LedAPI {
	return &ledAPI{
		host: host,
	}
}

func (la *ledAPI) getClient() led_service.LedServiceClient {
	conn, err := grpc.Dial(la.host, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return led_service.NewLedServiceClient(conn)
}

func (la *ledAPI) LightLED() error {
	laClient := la.getClient()

	_, err := laClient.SwitchLED(context.Background(), &emptypb.Empty{})

	return err
}

func (la *ledAPI) SwitchLED() error {
	laClient := la.getClient()

	_, err := laClient.SwitchLED(context.Background(), &emptypb.Empty{})

	return err
}
