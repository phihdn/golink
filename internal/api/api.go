package api

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger/v2"

	_ "github.com/phihdn/golink/docs" // registers the generated Swagger spec
	"github.com/phihdn/golink/internal/config"
	"github.com/phihdn/golink/internal/handler"
	"github.com/phihdn/golink/internal/repository"
	"github.com/phihdn/golink/internal/service"
	"github.com/phihdn/golink/pkg/validator"
)

// Engine is the interface of engine struct
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type engine struct {
	app         *echo.Echo
	cfg         *config.Config
	redisClient *redis.Client
}

// NewEngine create API Engine
func NewEngine(cfg *config.Config, redisClient *redis.Client) Engine {
	e := &engine{app: echo.New(), cfg: cfg, redisClient: redisClient}
	e.app.Validator = validator.New()
	e.initRoutes()
	return e
}

// Start starts the HTTP server on the configured port
func (e *engine) Start() error { return e.app.Start(":" + e.cfg.Port) }

func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.app.ServeHTTP(w, req)
}

func (e *engine) initRoutes() {
	// khoi tao handler
	healthCheckRepo := repository.NewHealthCheck(e.redisClient)
	healthSvc := service.NewHealth(e.cfg.ServiceName, e.cfg.InstanceID, healthCheckRepo)
	healthHandler := handler.NewHealth(healthSvc)

	genPassSvc := service.NewGenPass()
	genPassHandler := handler.NewGenPass(genPassSvc)

	urlStorageRepo := repository.NewURLStorage(e.redisClient)
	shortenURLSvc := service.NewShortenUrl(urlStorageRepo, genPassSvc)
	shortenHandler := handler.NewShortenLink(shortenURLSvc)

	e.app.GET("/health-check", healthHandler.Check)
	e.app.GET("/swagger/*", echoSwagger.EchoWrapHandler())
	e.app.GET("/v1/gen-pass", genPassHandler.GeneratePassword)
	e.app.POST("/v1/links/shorten", shortenHandler.CreateShortenLink)
}
