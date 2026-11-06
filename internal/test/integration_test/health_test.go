package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/phihdn/golink/internal/api"
	"github.com/phihdn/golink/internal/config"
)

func TestHealthCheck(t *testing.T) {
	srv := api.New(config.Config{ServiceName: "golink", InstanceID: "abc-123"})

	req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"OK","service_name":"golink","instance_id":"abc-123"}`, rec.Body.String())
}
