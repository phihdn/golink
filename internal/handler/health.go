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
//	@Description	Returns the service name and instance id. Returns 500 when Redis does not answer a ping.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	model.Health
//	@Failure		500	{object}	map[string]string
//	@Router			/health-check [get]
func (h *Health) Check(c *echo.Context) error {
	res, err := h.svc.Check(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	return c.JSON(http.StatusOK, res)
}
