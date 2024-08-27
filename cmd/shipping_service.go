package main

import (
	"flag"
	"log"
	"shipping/internal/infra/config"
	"shipping/internal/infra/http"
)

func main() {
	configPath := flag.String("config", "./config/core.yaml", "load configurations from a file")
	flag.Parse()

	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	http.RunServer(config)
}
