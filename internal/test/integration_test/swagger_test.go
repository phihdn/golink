package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/phihdn/golink/internal/api"
	"github.com/phihdn/golink/internal/config"
)

func TestSwaggerDoc(t *testing.T) {
	srv := api.New(config.Config{})

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "/health-check")
	assert.Contains(t, rec.Body.String(), "/v1/gen-pass")
}
