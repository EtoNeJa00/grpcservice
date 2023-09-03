package app

import (
	"context"
	"fmt"
	"log"

	"GRPCService/internal/app/repository"
	"GRPCService/internal/app/usecase"
	"GRPCService/internal/app/utilities/prom_metrics"
	"GRPCService/internal/config"
	"GRPCService/internal/endpoints"
	"GRPCService/internal/transport/grpc"
	"GRPCService/internal/transport/http_transport"

	"github.com/prometheus/client_golang/prometheus"
)

func StartApp(ctx context.Context, conf config.Config) (func(), error) {
	var servers []grpc.GRPCServer

	reg := prometheus.NewRegistry()
	m, err := prom_metrics.CreateMetrics(reg)

	serverMC, err := startMCApp(ctx, conf, m)
	if err != nil {
		return nil, err
	}

	servers = append(servers, serverMC)

	serverIS, err := startISApp(ctx, conf, m)
	if err != nil {
		return nil, err
	}

	http_transport.Start(reg, conf)

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

func startISApp(ctx context.Context, conf config.Config, m *prom_metrics.Metrics) (grpc.GRPCServer, error) {
	r := repository.NewInnerStorageRepository()
	uc := usecase.NewUsecase(r)

	PromUC, err := prom_metrics.NewPrometheusMiddleware(uc, m, "internal_storage")
	if err != nil {
		return nil, err
	}

	enp := endpoints.NewEndpoint(ctx, PromUC)

	sr, err := grpc.NewGRPCServer(conf.PortIS, enp)

	log.Printf("start grpc internal storage server on: " + conf.PortIS)

	return sr, err
}

func startMCApp(ctx context.Context, conf config.Config, m *prom_metrics.Metrics) (grpc.GRPCServer, error) {
	r, err := repository.NewMemcacheRepository(conf.MCServerAddr)
	if err != nil {
		return nil, err
	}

	uc := usecase.NewUsecase(r)

	PromUC, err := prom_metrics.NewPrometheusMiddleware(uc, m, "memcached_storage")
	if err != nil {
		return nil, err
	}

	enp := endpoints.NewEndpoint(ctx, PromUC)
	sr, err := grpc.NewGRPCServer(conf.PortMC, enp)

	log.Printf("start grpc memcached server on: " + conf.PortMC)

	return sr, err
}
