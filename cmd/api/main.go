package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/phihdn/golink/internal/api"
	"github.com/phihdn/golink/internal/config"
	"github.com/phihdn/golink/pkg/redis"
)

//	@title			golink API
//	@version		1.0
//	@description	URL shortener service.
//	@BasePath		/

func main() {
	// Handlers log through zerolog.Ctx(request context). Without a logger in
	// the context, fall back to the global one, which writes JSON to stderr.
	zerolog.DefaultContextLogger = &log.Logger

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
