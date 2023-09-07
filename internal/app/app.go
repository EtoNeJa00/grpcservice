package app

import (
	"context"
	"fmt"
	"log"

	"GRPCService/config"
	"GRPCService/internal/app/repository"
	"GRPCService/internal/app/usecase"
	"GRPCService/internal/app/utilities/prommetrics"
	"GRPCService/internal/endpoints"
	"GRPCService/internal/transport/grpc"
	"GRPCService/internal/transport/httptransport"

	"github.com/prometheus/client_golang/prometheus"
)

func StartApp(ctx context.Context, conf *config.Config) (func(), error) {
	var servers []grpc.GRPCServer

	reg := prometheus.NewRegistry()

	m, err := prommetrics.CreateMetrics(reg)
	if err != nil {
		return nil, err
	}

	serverMC, err := startMCApp(ctx, conf, m)
	if err != nil {
		return nil, err
	}

	servers = append(servers, serverMC)

	serverIS, err := startISApp(ctx, conf, m)
	if err != nil {
		return nil, err
	}

	servers = append(servers, serverIS)

	serverSc, err := startScApp(ctx, conf, m)
	if err != nil {
		return nil, err
	}

	servers = append(servers, serverSc)

	httptransport.Start(reg, conf)

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
			s.StopServer()
		}
	}, nil
}

func startISApp(ctx context.Context, conf *config.Config, m *prommetrics.Metrics) (grpc.GRPCServer, error) {
	r := repository.NewInnerStorageRepository()
	uc := usecase.NewUsecase(r)

	PromUC, err := prommetrics.NewPrometheusMiddleware(uc, m, "internal_storage")
	if err != nil {
		return nil, err
	}

	enp := endpoints.NewEndpoint(ctx, PromUC)

	sr, err := grpc.NewGRPCServer(conf.PortIS, enp)

	log.Printf("start grpc internal storage server on: " + conf.PortIS)

	return sr, err
}

func startMCApp(ctx context.Context, conf *config.Config, m *prommetrics.Metrics) (grpc.GRPCServer, error) {
	r, err := repository.NewMemcacheRepository(conf.MCServerAddr)
	if err != nil {
		return nil, err
	}

	uc := usecase.NewUsecase(r)

	PromUC, err := prommetrics.NewPrometheusMiddleware(uc, m, "memcached_storage")
	if err != nil {
		return nil, err
	}

	enp := endpoints.NewEndpoint(ctx, PromUC)
	sr, err := grpc.NewGRPCServer(conf.PortMC, enp)

	log.Printf("start grpc memcached server on: " + conf.PortMC)

	return sr, err
}

func startScApp(ctx context.Context, conf *config.Config, m *prommetrics.Metrics) (grpc.GRPCServer, error) {
	r, err := repository.NewScyllaRepository(ctx, conf.ScyllaAddr)
	if err != nil {
		return nil, err
	}

	uc := usecase.NewUsecase(r)

	PromUC, err := prommetrics.NewPrometheusMiddleware(uc, m, "scylla_storage")
	if err != nil {
		return nil, err
	}

	enp := endpoints.NewEndpoint(ctx, PromUC)
	sr, err := grpc.NewGRPCServer(conf.PortSc, enp)

	log.Printf("start grpc scylla server on: " + conf.PortSc)

	return sr, err
}
