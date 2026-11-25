// Package validator plugs go-playground/validator into Echo. Gin ships it built
// in (tag `binding`); Echo only defines the echo.Validator interface, so the app
// registers this adapter and handlers call c.Validate after c.Bind.
package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type customValidator struct {
	v *validator.Validate
}

// New returns an echo.Validator that checks `validate` struct tags.
func New() echo.Validator {
	return &customValidator{v: validator.New(validator.WithRequiredStructEnabled())}
}

func (cv *customValidator) Validate(i any) error {
	return cv.v.Struct(i)
}
