package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PortIS string `envconfig:"INTERNAL_STORAGE_GRPC_PORT" default:":5300"`
	PortMC string `envconfig:"MEMCACHE_GRPC_PORT" default:":8080"`

	PortPrometheus string `envconfig:"PROMETHEUS_PORT" default:":8888"`

	MCServerAddr string `envconfig:"MEMCACHED_ADDR" default:"localhost:11211"`
}

func GetConfig() (Config, error) {
	conf := Config{}

	err := envconfig.Process("", &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}
