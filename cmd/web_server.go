package main

import (
	"flag"
	"log"
	"shipping/internal/app/route"
	"shipping/internal/infra/config"
)

func main() {
	configPath := flag.String("config", "./config/core.yaml", "load configurations from a file")
	flag.Parse()

	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	r := route.SetupRouter(config)
	r.Run(config.App.ListenAddr)
}
