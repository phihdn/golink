package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"

	"github.com/phihdn/golink/internal/service/mocks"
)

func TestShortenLink_CreateShortenLink(t *testing.T) {
	t.Parallel()

	const (
		invalidInput = `{"error":"Invalid Input"}`
		successBody  = `{"code":"abc1234","message":"Shorten URL generated successfully!"}`
	)

	testCases := []struct {
		name         string
		body         string
		setup        func(ctx context.Context, svc *mocks.ShortenUrl)
		expectedCode int
		expectedBody string
	}{
		{
			name: "success with exp",
			body: `{"exp":604800,"url":"https://example.com"}`,
			setup: func(ctx context.Context, svc *mocks.ShortenUrl) {
				svc.EXPECT().
					CreateShortenLink(ctx, "https://example.com", 604800*time.Second).
					Return("abc1234", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: successBody,
		},
		{
			name: "success without exp",
			body: `{"url":"https://example.com"}`,
			setup: func(ctx context.Context, svc *mocks.ShortenUrl) {
				svc.EXPECT().
					CreateShortenLink(ctx, "https://example.com", time.Duration(0)).
					Return("abc1234", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: successBody,
		},
		{
			name:         "malformed body",
			body:         `{"url":`,
			expectedCode: http.StatusBadRequest,
			expectedBody: invalidInput,
		},
		{
			name:         "missing url",
			body:         `{"exp":10}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: invalidInput,
		},
		{
			name:         "relative url",
			body:         `{"url":"/just/a/path"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: invalidInput,
		},
		{
			name:         "unsupported scheme",
			body:         `{"url":"ftp://example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: invalidInput,
		},
		{
			name:         "missing host",
			body:         `{"url":"https://"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: invalidInput,
		},
		{
			name:         "port without host",
			body:         `{"url":"https://:80"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: invalidInput,
		},
		{
			name:         "negative exp",
			body:         `{"exp":-1,"url":"https://example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: invalidInput,
		},
		{
			name:         "exp too large",
			body:         `{"exp":9223372037,"url":"https://example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: invalidInput,
		},
		{
			name: "service error",
			body: `{"url":"https://example.com"}`,
			setup: func(ctx context.Context, svc *mocks.ShortenUrl) {
				svc.EXPECT().
					CreateShortenLink(ctx, "https://example.com", time.Duration(0)).
					Return("", errors.New("redis down"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			svc := mocks.NewShortenUrl(t)
			if tc.setup != nil {
				tc.setup(ctx, svc)
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(tc.body)).
				WithContext(ctx)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			err := NewShortenLink(svc).CreateShortenLink(c)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
		})
	}
}
