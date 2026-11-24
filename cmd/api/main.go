package main

import (
	"github.com/phihdn/golink/internal/api"
	"github.com/phihdn/golink/internal/config"
	"github.com/phihdn/golink/pkg/redis"
)

//	@title			golink API
//	@version		1.0
//	@description	URL shortener service.
//	@BasePath		/

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}

	app := api.NewEngine(&cfg, redisClient)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
