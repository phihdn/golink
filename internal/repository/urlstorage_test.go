package repository

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	redisPkg "github.com/phihdn/golink/pkg/redis"
)

// initClosedRedis returns a client that fails every command, to test error paths.
func initClosedRedis(t *testing.T) *redis.Client {
	r := redisPkg.InitMockRedis(t)
	_ = r.Close()
	return r
}

func TestUrlStorage_StoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		exp       time.Duration
		setupMock func(ctx context.Context) *redis.Client

		expectedErr error
		verifyFunc  func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",
			exp:  time.Hour,
			setupMock: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},

			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				res, err := r.Get(ctx, "abcxyz").Result()
				assert.NoError(t, err)
				assert.Equal(t, "https://google.com", res)
				assert.Equal(t, time.Hour, r.TTL(ctx, "abcxyz").Val())
			},
		},
		{
			name: "no expiration",
			setupMock: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},

			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				// redis reports -1 for a key without expiration
				assert.Equal(t, time.Duration(-1), r.TTL(ctx, "abcxyz").Val())
			},
		},
		{
			name: "key existed",
			exp:  time.Hour,
			setupMock: func(ctx context.Context) *redis.Client {
				r := redisPkg.InitMockRedis(t)
				r.Set(ctx, "abcxyz", "https://old.com", 0)
				return r
			},

			expectedErr: ErrKeyExisted,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				res, err := r.Get(ctx, "abcxyz").Result()
				assert.NoError(t, err)
				assert.Equal(t, "https://old.com", res, "existing url must not be overwritten")
			},
		},
		{
			name: "redis error",
			setupMock: func(ctx context.Context) *redis.Client {
				return initClosedRedis(t)
			},

			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mock := tc.setupMock(ctx)
			testRepo := NewURLStorage(mock)
			err := testRepo.StoreURL(ctx, "abcxyz", "https://google.com", tc.exp)
			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mock)
			}
		})
	}
}

func TestUrlStorage_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context) *redis.Client

		expectedURL string
		expectedErr error
	}{
		{
			name: "normal case",
			setupMock: func(ctx context.Context) *redis.Client {
				r := redisPkg.InitMockRedis(t)
				r.Set(ctx, "abcxyz", "https://google.com", 0)
				return r
			},
			expectedURL: "https://google.com",
		},
		{
			name: "not found",
			setupMock: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			expectedErr: ErrKeyNotFound,
		},
		{
			name: "redis error",
			setupMock: func(ctx context.Context) *redis.Client {
				return initClosedRedis(t)
			},
			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			testRepo := NewURLStorage(tc.setupMock(ctx))
			url, err := testRepo.GetURL(ctx, "abcxyz")
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedURL, url)
		})
	}
}
