package grpc

import (
	"net"

	"GRPCService/internal/endpoints"
	"GRPCService/internal/transport/grpc/generated/record"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer interface {
	StartServer() error
	StopServer()
}

type gRPCServer struct {
	enp        endpoints.GrpcEnp
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

	var opts []grpc.ServerOption
	s.grpcServer = grpc.NewServer(opts...)

	record.RegisterRecordsServer(s.grpcServer, s.enp)
	reflection.Register(s.grpcServer)

	return s.grpcServer.Serve(listener)
}

func (s *gRPCServer) StopServer() {
	s.grpcServer.GracefulStop()
}
