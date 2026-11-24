package handler

import (
	"math"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"

	"github.com/phihdn/golink/internal/service"
)

const (
	shortenSuccessMessage = "Shorten URL generated successfully!"
	// maxExpSeconds is the largest exp that still fits in a time.Duration.
	maxExpSeconds = math.MaxInt64 / int64(time.Second)
)

type ShortenLink interface {
	CreateShortenLink(c *echo.Context) error
}

type shortenLink struct {
	svc service.ShortenUrl
}

func NewShortenLink(svc service.ShortenUrl) ShortenLink {
	return &shortenLink{svc: svc}
}

type shortenLinkInput struct {
	// Exp is the link lifetime in seconds. 0 or omitted means the link never expires.
	Exp int64  `json:"exp" example:"604800"`
	URL string `json:"url" example:"https://example.com/some/long/path" validate:"required,url"`
}

type shortenLinkResponse struct {
	Code    string `json:"code"    example:"aB3dE6g"`
	Message string `json:"message" example:"Shorten URL generated successfully!"`
}

// CreateShortenLink   Generate shorten link
// @Summary      Generate shorten url based on original url
// @Description  Generate a unique 7-character code for the url. exp is the lifetime in seconds, 0 or omitted means no expiration.
// @Tags         link
// @Accept       application/json
// @Produce      application/json
// @Param        input body shortenLinkInput true "Input required"
// @Success      200 {object} shortenLinkResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /v1/links/shorten [post]
func (h *shortenLink) CreateShortenLink(c *echo.Context) error {
	// lay input
	input := &shortenLinkInput{}
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Input"})
	}
	if err := c.Validate(input); err != nil || input.Exp < 0 || input.Exp > maxExpSeconds {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Input"})
	}

	// goi svc
	// *echo.Context is not a context.Context (unlike *gin.Context), so pass the
	// request context. It is the one context every layer below must receive.
	code, err := h.svc.CreateShortenLink(
		c.Request().Context(),
		input.URL,
		time.Duration(input.Exp)*time.Second,
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": "Internal Server Error"},
		)
	}

	// return resp
	return c.JSON(http.StatusOK, shortenLinkResponse{Code: code, Message: shortenSuccessMessage})
}
