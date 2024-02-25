package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	zigbee_service_pb "github.com/vadim-dmitriev/mirror-zigbee-api/pkg/zigbee-service"
)

// HTTPServer is http server struct
type HTTPServer struct {
	httpPort int
	grpcPort int

	mux *runtime.ServeMux
}

// NewHTTP creates new HTTPServer object.
func NewHTTP(httpPort int, grpcPort int) (*HTTPServer, error) {
	mux := runtime.NewServeMux()

	s := &HTTPServer{
		httpPort: httpPort,
		grpcPort: grpcPort,
		mux:      mux,
	}

	return s, nil
}

// Run starts server
func (s *HTTPServer) Run(ctx context.Context) error {
	defer log.Println("HTTP server stopped")

	// TODO: get host and port from config
	httpAddr := fmt.Sprintf(":%d", s.httpPort)
	listener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		return fmt.Errorf("failed to create tcp listener on '%s' addr: %w", httpAddr, err)
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := zigbee_service_pb.RegisterZigbeeServiceHandlerFromEndpoint(ctx, s.mux, fmt.Sprintf(":%d", s.grpcPort), opts); err != nil {
		return fmt.Errorf("failed to register reverse-proxy: %w", err)
	}

	log.Printf("HTTP server listening at %v", listener.Addr())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigCh

		log.Printf("got signal %v, attempting graceful shutdown HTTP server", sig)

		if err := listener.Close(); err != nil {
			os.Exit(1)
		}
	}()

	// Перед запуском настраиваем правила CORS.
	if err := http.Serve(listener, cors.AllowAll().Handler(s.mux)); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}

	return nil
}
