package api

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/phihdn/golink/internal/config"
	"github.com/phihdn/golink/internal/handler"
	"github.com/phihdn/golink/internal/service"
)

type Engine interface {
	http.Handler
	Start(address string) error
}

type api struct {
	*echo.Echo
	cfg config.Config
}

func New(cfg config.Config) Engine {
	a := &api{Echo: echo.New(), cfg: cfg}
	a.registerRoutes()
	return a
}

func (a *api) registerRoutes() {
	healthHandler := handler.NewHealth(service.NewHealth(a.cfg.ServiceName, a.cfg.InstanceID))
	a.GET("/health-check", healthHandler.Check)
}

func (a *api) Start(address string) error {
	if address == "" {
		address = ":" + a.cfg.Port
	}
	return a.Echo.Start(address)
}
