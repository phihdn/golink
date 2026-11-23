package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/phihdn/golink/internal/model"
	"github.com/phihdn/golink/internal/service/mocks"
)

func TestHealth_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		health       model.Health
		err          error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "healthy",
			health:       model.Health{Message: "OK", ServiceName: "golink", InstanceID: "abc-123"},
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"OK","service_name":"golink","instance_id":"abc-123"}`,
		},
		{
			name:         "redis down",
			err:          errors.New("connection refused"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewHealth(t)
			svc.EXPECT().Check(mock.Anything).Return(tc.health, tc.err)

			req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			require.NoError(t, NewHealth(svc).Check(c))

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
		})
	}
}
