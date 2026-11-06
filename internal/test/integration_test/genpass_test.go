package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/phihdn/golink/internal/api"
	"github.com/phihdn/golink/internal/config"
)

func TestGenPass(t *testing.T) {
	srv := api.New(config.Config{})

	req := httptest.NewRequest(http.MethodGet, "/v1/gen-pass", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Regexp(t, `^[a-zA-Z0-9]{10}$`, rec.Body.String())
}

func TestGenPass_UnversionedPathNotFound(t *testing.T) {
	srv := api.New(config.Config{})

	req := httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
