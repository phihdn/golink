package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/phihdn/golink/internal/service"
)

const passwordLength = 10

type GenPass struct {
	svc service.GenPass
}

func NewGenPass(svc service.GenPass) *GenPass {
	return &GenPass{svc: svc}
}

// GeneratePassword godoc
//
//	@Summary		Generate a password
//	@Description	Returns a random 10-character alphanumeric password as plain text.
//	@Tags			password
//	@Produce		plain
//	@Success		200	{string}	string	"aB3dE6gH9j"
//	@Failure		500	{object}	model.ErrorResponse
//	@Router			/v1/gen-pass [get]
func (g *GenPass) GeneratePassword(c *echo.Context) error {
	password, err := g.svc.GeneratePassword(passwordLength)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, password)
}
