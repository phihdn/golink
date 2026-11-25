package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"

	"github.com/phihdn/golink/internal/service/mocks"
)

func TestGenPass_GeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		password     string
		serviceErr   error
		expectedCode int
		expectedBody string
		expectedErr  error
		expectedLog  string
	}{
		{
			name:         "success",
			password:     "aB3dE6gH9j",
			expectedCode: http.StatusOK,
			expectedBody: "aB3dE6gH9j",
		},
		{
			name:        "service error",
			serviceErr:  errors.New("rand failure"),
			expectedErr: errors.New("rand failure"),
			expectedLog: "generate password failed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewGenPass(t)
			svc.EXPECT().GeneratePassword(passwordLength).Return(tc.password, tc.serviceErr)

			ctx, logs := newLoggedContext(t)
			req := httptest.NewRequest(http.MethodGet, "/v1/gen-pass", nil).WithContext(ctx)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)

			err := NewGenPass(svc).GeneratePassword(c)

			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedErr == nil {
				assert.Equal(t, tc.expectedCode, rec.Code)
				assert.Equal(t, tc.expectedBody, rec.Body.String())
			}
			assertErrorLog(t, logs, tc.expectedLog, "rand failure")
		})
	}
}
