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

func (g *GenPass) GeneratePassword(c *echo.Context) error {
	password, err := g.svc.GeneratePassword(passwordLength)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, password)
}
