package zigbee_service

/*

grpcurl -plaintext \
	-d '{"name":"lamp-1", "states": [{"name": "state", "type":3, "text_value": "ON"}, {"name": "brightness", "type":2, "number_value": 2500}]}' \
	localhost:9080 zigbee_service.ZigbeeService/SetDeviceState

*/

import (
	"context"
	"fmt"
	"log"

	"github.com/vadim-dmitriev/mirror-zigbee-api/internal/domain"
	zigbee_service_pb "github.com/vadim-dmitriev/mirror-zigbee-api/pkg/zigbee-service"
)

func (zs *zigbeeService) SetDeviceState(ctx context.Context, request *zigbee_service_pb.SetDeviceStateRequest) (*zigbee_service_pb.SetDeviceStateResponse, error) {
	log.Printf("'SetDeviceState' executed with request '%s'\n", request)

	states := make([]*domain.DeviceState, 0, len(request.GetStates()))
	for _, stateRaw := range request.GetStates() {
		state, err := mapState(stateRaw)
		if err != nil {
			return nil, fmt.Errorf("failed to map device state with name '%s': %w", stateRaw.GetName(), err)
		}

		states = append(states, state)
	}

	_, err := zs.zigbeeClient.SetDeviceState(ctx, request.GetName(), states)
	if err != nil {
		return nil, fmt.Errorf("failed to set device state: %w", err)
	}

	return &zigbee_service_pb.SetDeviceStateResponse{}, nil
}

func mapState(statePb *zigbee_service_pb.State) (*domain.DeviceState, error) {
	var stateValue interface{}

	switch statePb.GetType() {
	case zigbee_service_pb.State_Type_UNKNOWN:
		return nil, fmt.Errorf("unknown state type")
	case zigbee_service_pb.State_Type_BOOLEAN:
		stateValue = statePb.GetBooleanValue().GetValue()
	case zigbee_service_pb.State_Type_NUMBER:
		stateValue = statePb.GetNumberValue().GetValue()
	case zigbee_service_pb.State_Type_TEXT:
		stateValue = statePb.GetTextValue().GetValue()
	default:
		return nil, fmt.Errorf("unexpected state type")
	}

	state := &domain.DeviceState{
		Name:  statePb.GetName(),
		Value: stateValue,
	}

	return state, nil
}

func mapStatePb(state *domain.DeviceState) *zigbee_service_pb.State {
	statePbType := zigbee_service_pb.State_Type_UNKNOWN

	switch state.Type {
	case domain.DeviceStateType_BOOLEAN:
		statePbType = zigbee_service_pb.State_Type_BOOLEAN
	case domain.DeviceStateType_INTEGER:
		statePbType = zigbee_service_pb.State_Type_NUMBER
	case domain.DeviceStateType_STRING:
		statePbType = zigbee_service_pb.State_Type_TEXT
	}

	statePb := &zigbee_service_pb.State{
		Type: statePbType,
		Name: state.Name,
		// Value: state.Value,
	}

	return statePb
}
