package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog"

	"github.com/phihdn/golink/internal/model"
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
//	@Failure		500	{object}	model.InternalServerError
//	@Router			/health-check [get]
func (h *Health) Check(c *echo.Context) error {
	// *echo.Context is not a context.Context (unlike *gin.Context), so pass the
	// request context. It is the one context every layer below must receive.
	res, err := h.svc.Check(c.Request().Context())
	if err != nil {
		zerolog.Ctx(c.Request().Context()).Error().Err(err).Msg("health check failed")
		return c.JSON(
			http.StatusInternalServerError,
			model.HTTPError{Error: "Internal Server Error"},
		)
	}

	return c.JSON(http.StatusOK, res)
}
