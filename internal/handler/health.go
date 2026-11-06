package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/phihdn/golink/internal/service"
)

type Health struct {
	svc service.Health
}

func NewHealth(svc service.Health) *Health {
	return &Health{svc: svc}
}

func (h *Health) Check(c *echo.Context) error {
	return c.JSON(http.StatusOK, h.svc.Check())
}
