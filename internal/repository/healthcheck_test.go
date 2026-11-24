package repository

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	redisPkg "github.com/phihdn/golink/pkg/redis"
)

func TestHealthCheck_Ping(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupMock   func() *redis.Client
		expectedErr error
	}{
		{
			name:      "redis is up",
			setupMock: func() *redis.Client { return redisPkg.InitMockRedis(t) },
		},
		{
			name:        "redis is down",
			setupMock:   func() *redis.Client { return initClosedRedis(t) },
			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := NewHealthCheck(tc.setupMock()).Ping(t.Context())
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
