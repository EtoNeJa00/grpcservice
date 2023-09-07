package main

import (
	"context"
	"log"

	"GRPCService/config"
	"GRPCService/internal/app"
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

	<-ctx.Done()

	log.Print("cancelled")

	stop()
}
