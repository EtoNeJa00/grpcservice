package http_transport

import (
	"log"
	"net/http"

	"GRPCService/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(reg *prometheus.Registry, config config.Config) {
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	go func() {
		err := http.ListenAndServe(config.PortPrometheus, nil)
		if err != nil {
			log.Println(err)
		}
	}()

	return
}
