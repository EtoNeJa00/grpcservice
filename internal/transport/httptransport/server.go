package httptransport

import (
	"log"
	"net/http"

	"GRPCService/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(reg *prometheus.Registry, conf *config.Config) {
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	go func() {
		err := http.ListenAndServe(conf.PortPrometheus, nil)
		if err != nil {
			log.Println(err)
		}
	}()
}
