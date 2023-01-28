package app

import (
	"context"
	"fmt"

	"github.com/EtoNeJa00/GRPCService/internal/app/repository"
	"github.com/EtoNeJa00/GRPCService/internal/app/usecase"
	"github.com/EtoNeJa00/GRPCService/internal/config"
	"github.com/EtoNeJa00/GRPCService/internal/endpoints"
	"github.com/EtoNeJa00/GRPCService/internal/transport/grpc"
)

func StartApp(ctx context.Context, conf config.Config) (func(), error) {
	var servers []grpc.GRPCServer

	serverMC, err := startMCApp(ctx, conf)
	if err != nil {
		return nil, err
	}

	servers = append(servers, serverMC)

	serverIS, err := startISApp(ctx, conf)
	if err != nil {
		return nil, err
	}

	servers = append(servers, serverIS)

	for _, s := range servers {
		s := s
		go func() {
			err := s.StartServer()
			if err != nil {
				fmt.Print(err)
			}
		}()
	}

	return func() {
		for _, s := range servers {
			(s).StopServer()
		}
	}, nil
}

func startISApp(ctx context.Context, conf config.Config) (grpc.GRPCServer, error) {
	r := repository.NewInnerStorageRepository(ctx)
	uc := usecase.NewUsecase(r)
	enp := endpoints.NewEndpoint(ctx, uc)

	sr, err := grpc.NewGRPCServer(conf.PortIS, enp)

	return sr, err
}

func startMCApp(ctx context.Context, conf config.Config) (grpc.GRPCServer, error) {
	r, err := repository.NewMemcacheRepository(conf.MCServerAddr)
	if err != nil {
		return nil, err
	}

	uc := usecase.NewUsecase(r)
	enp := endpoints.NewEndpoint(ctx, uc)
	sr, err := grpc.NewGRPCServer(conf.PortMC, enp)

	return sr, err
}
