package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog"

	"github.com/phihdn/golink/internal/model"
	"github.com/phihdn/golink/internal/service"
)

const shortenSuccessMessage = "Shorten URL generated successfully!"

type ShortenLink interface {
	CreateShortenLink(c *echo.Context) error
	Redirect(c *echo.Context) error
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
		zerolog.Ctx(c.Request().Context()).Error().
			Err(err).
			Str("url", input.URL).
			Msg("create shorten link failed")
		return c.JSON(
			http.StatusInternalServerError,
			model.HTTPError{Error: "Internal Server Error"},
		)
	}

	// return resp
	return c.JSON(http.StatusOK, shortenLinkResponse{Code: code, Message: shortenSuccessMessage})
}

// Redirect   Redirect to the original url
// @Summary      Redirect a short code to its original url
// @Description  Look up the code and answer 302 Found with the original url in the Location header.
// @Tags         link
// @Produce      application/json
// @Param        code path string true "Short code" example(aB3dE6g)
// @Success      302
// @Failure      400 {object} model.BadRequestError
// @Failure      404 {object} model.NotFoundError
// @Failure      500 {object} model.InternalServerError
// @Router       /v1/links/redirect/{code} [get]
func (h *shortenLink) Redirect(c *echo.Context) error {
	// lay input
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, model.HTTPError{Error: "Invalid Input"})
	}

	// goi svc
	url, err := h.svc.GetUrlFromCode(c.Request().Context(), code)
	if errors.Is(err, service.ErrLinkNotFound) {
		return c.JSON(http.StatusNotFound, model.HTTPError{Error: "Not Found"})
	}
	if err != nil {
		zerolog.Ctx(c.Request().Context()).Error().
			Err(err).
			Str("code", code).
			Msg("get url from code failed")
		return c.JSON(
			http.StatusInternalServerError,
			model.HTTPError{Error: "Internal Server Error"},
		)
	}

	// redirect
	return c.Redirect(http.StatusFound, url)
}
