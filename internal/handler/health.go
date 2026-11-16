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

// Check godoc
//
//	@Summary		Health check
//	@Description	Returns the service name and instance id.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	model.Health
//	@Router			/health-check [get]
func (h *Health) Check(c *echo.Context) error {
	return c.JSON(http.StatusOK, h.svc.Check())
}
