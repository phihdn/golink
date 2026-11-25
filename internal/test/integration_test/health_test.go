package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/phihdn/golink/internal/config"
)

func healthCheck(srv http.Handler) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

func TestHealthCheck(t *testing.T) {
	srv, _ := newServer(t, config.Config{ServiceName: "golink", InstanceID: "abc-123"})

	rec := healthCheck(srv)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(
		t,
		`{"message":"OK","service_name":"golink","instance_id":"abc-123"}`,
		rec.Body.String(),
	)
}

func TestHealthCheck_RedisDown(t *testing.T) {
	srv, mr := newServer(t, config.Config{ServiceName: "golink", InstanceID: "abc-123"})
	mr.Close()

	rec := healthCheck(srv)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"error":"Internal Server Error"}`, rec.Body.String())
}
