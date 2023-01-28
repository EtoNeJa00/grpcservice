package main

import (
	"context"
	"log"

	"github.com/EtoNeJa00/GRPCService/internal/app"
	"github.com/EtoNeJa00/GRPCService/internal/config"
)

func main() {
	ctx := context.Background()

	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	stop, err := app.StartApp(ctx, conf)
	if err != nil {
		log.Fatal(err)
	}

	select {
	case <-ctx.Done():
		log.Print("cancelled")
	}

	stop()
}
