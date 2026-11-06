package main

import (
	"log"

	"github.com/phihdn/golink/internal/api"
	"github.com/phihdn/golink/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := api.New(cfg).Start(""); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
