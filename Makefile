PROTO_PATH=./api

.generate:
	protoc -I=$(PROTO_PATH) \
		--go_out=./pkg --go_opt paths=source_relative \
		--go-grpc_out=./pkg --go-grpc_opt paths=source_relative \
		--grpc-gateway_out=./pkg --grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		zigbee-service/zigbee-service.proto

generate: .generate

.run:
	go run ./cmd

run: .run
