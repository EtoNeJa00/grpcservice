package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/EtoNeJa00/GRPCService/internal/endpoints"
	"github.com/EtoNeJa00/GRPCService/internal/transport/grpc/generated/record"
)

type GRPCServer interface {
	StartServer() error
	StopServer()
}

type gRPCServer struct {
	enp        endpoints.GrpcEnp
	listener   net.Listener
	grpcServer *grpc.Server
	port       string
}

func NewGRPCServer(port string, enp endpoints.GrpcEnp) (GRPCServer, error) {
	return &gRPCServer{
		port: port,
		enp:  enp,
	}, nil
}

func (s *gRPCServer) StartServer() (err error) {
	listener, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	log.Printf("service start on:%s\n", listener.Addr())

	opts := []grpc.ServerOption{}
	s.grpcServer = grpc.NewServer(opts...)

	record.RegisterRecordsServer(s.grpcServer, s.enp)

	return s.grpcServer.Serve(listener)
}

func (s *gRPCServer) StopServer() {
	s.grpcServer.GracefulStop()
}
