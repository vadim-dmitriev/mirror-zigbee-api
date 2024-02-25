package server

// grpcurl -plaintext localhost:8080 list

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var _ grpc.ServiceRegistrar = &GRPCServer{}
var _ reflection.GRPCServer = &GRPCServer{}

// GRPCServer is grpc server struct
type GRPCServer struct {
	listener net.Listener
	grpc     *grpc.Server
}

// NewGRPC creates new GRPCServer object.
func NewGRPC(port int) (*GRPCServer, error) {
	addr := fmt.Sprintf(":%d", port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create tcp listener on '%s' addr: %w", addr, err)
	}

	s := &GRPCServer{
		listener: listener,
		grpc:     grpc.NewServer(),
	}

	reflection.Register(s)

	return s, nil
}

// Run starts server
func (s *GRPCServer) Run() {
	defer log.Println("server stopped")

	log.Printf("GRPC server listening at %v", s.listener.Addr())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigCh

		log.Printf("got signal %v, attempting graceful shutdown GRPC server", sig)

		s.grpc.GracefulStop()
		if err := s.listener.Close(); err != nil {
			os.Exit(1)
		}
	}()

	if err := s.grpc.Serve(s.listener); err != nil {
		log.Fatalf("failed to serve GRPC: %v", err)
	}
}

func (s *GRPCServer) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	s.grpc.RegisterService(sd, ss)
}

func (s *GRPCServer) GetServiceInfo() map[string]grpc.ServiceInfo {
	return s.grpc.GetServiceInfo()
}
