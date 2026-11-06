package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/phihdn/golink/internal/model"
	"github.com/phihdn/golink/internal/service/mocks"
)

func TestHealth_Check(t *testing.T) {
	svc := mocks.NewHealth(t)
	svc.EXPECT().Check().Return(model.Health{
		Message:     "OK",
		ServiceName: "golink",
		InstanceID:  "abc-123",
	})

	req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	require.NoError(t, NewHealth(svc).Check(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message":"OK","service_name":"golink","instance_id":"abc-123"}`, rec.Body.String())
}
