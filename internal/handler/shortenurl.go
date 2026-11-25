package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"

	"github.com/phihdn/golink/internal/model"
	"github.com/phihdn/golink/internal/service"
)

const shortenSuccessMessage = "Shorten URL generated successfully!"

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
	// Exp is the link lifetime in seconds, at most 7 days (604800).
	// 0 or omitted means the link never expires.
	Exp int64  `json:"exp" validate:"gte=0,lte=604800" example:"604800"`
	URL string `json:"url" validate:"required,url"     example:"https://example.com/some/long/path"`
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
// @Failure      400 {object} model.BadRequestError
// @Failure      500 {object} model.InternalServerError
// @Router       /v1/links/shorten [post]
func (h *shortenLink) CreateShortenLink(c *echo.Context) error {
	// lay input
	input := &shortenLinkInput{}
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, model.HTTPError{Error: "Invalid Input"})
	}
	if err := c.Validate(input); err != nil {
		return c.JSON(http.StatusBadRequest, model.HTTPError{Error: "Invalid Input"})
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
			model.HTTPError{Error: "Internal Server Error"},
		)
	}

	// return resp
	return c.JSON(http.StatusOK, shortenLinkResponse{Code: code, Message: shortenSuccessMessage})
}
