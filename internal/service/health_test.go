package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/phihdn/golink/internal/model"
	repomocks "github.com/phihdn/golink/internal/repository/mocks"
)

func TestHealth_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		pingErr        error
		expectedHealth model.Health
		expectedErr    error
	}{
		{
			name:           "redis is up",
			expectedHealth: model.Health{Message: "OK", ServiceName: "golink", InstanceID: "abc-123"},
		},
		{
			name:        "redis is down",
			pingErr:     errors.New("connection refused"),
			expectedErr: errors.New("connection refused"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			repo := repomocks.NewHealthCheck(t)
			repo.EXPECT().Ping(ctx).Return(tc.pingErr).Once()

			got, err := NewHealth("golink", "abc-123", repo).Check(ctx)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedHealth, got)
		})
	}
}
