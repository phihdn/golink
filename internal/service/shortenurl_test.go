package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/phihdn/golink/internal/repository"
	repomocks "github.com/phihdn/golink/internal/repository/mocks"
	"github.com/phihdn/golink/internal/service/mocks"
)

func TestShortenUrl_CreateShortenLink(t *testing.T) {
	t.Parallel()

	const url = "https://example.com/some/long/path?q=1"

	testCases := []struct {
		name         string
		exp          time.Duration
		setup        func(ctx context.Context, gen *mocks.GenPass, storage *repomocks.URLStorage)
		expectedCode string
		expectedErr  error
	}{
		{
			name: "stored on first try",
			exp:  time.Hour,
			setup: func(ctx context.Context, gen *mocks.GenPass, storage *repomocks.URLStorage) {
				gen.EXPECT().GeneratePassword(DefaultCodeLength).Return("abc1234", nil).Once()
				storage.EXPECT().StoreURL(ctx, "abc1234", url, time.Hour).Return(nil).Once()
			},
			expectedCode: "abc1234",
		},
		{
			name: "stored without expiration",
			setup: func(ctx context.Context, gen *mocks.GenPass, storage *repomocks.URLStorage) {
				gen.EXPECT().GeneratePassword(DefaultCodeLength).Return("abc1234", nil).Once()
				storage.EXPECT().StoreURL(ctx, "abc1234", url, time.Duration(0)).Return(nil).Once()
			},
			expectedCode: "abc1234",
		},
		{
			name: "retries when code already exists",
			setup: func(ctx context.Context, gen *mocks.GenPass, storage *repomocks.URLStorage) {
				gen.EXPECT().GeneratePassword(DefaultCodeLength).Return("taken00", nil).Once()
				gen.EXPECT().GeneratePassword(DefaultCodeLength).Return("free000", nil).Once()
				storage.EXPECT().StoreURL(ctx, "taken00", url, time.Duration(0)).Return(repository.ErrKeyExisted).Once()
				storage.EXPECT().StoreURL(ctx, "free000", url, time.Duration(0)).Return(nil).Once()
			},
			expectedCode: "free000",
		},
		{
			name: "retries when storage wraps the existed error",
			setup: func(ctx context.Context, gen *mocks.GenPass, storage *repomocks.URLStorage) {
				gen.EXPECT().GeneratePassword(DefaultCodeLength).Return("taken00", nil).Once()
				gen.EXPECT().GeneratePassword(DefaultCodeLength).Return("free000", nil).Once()
				storage.EXPECT().StoreURL(ctx, "taken00", url, time.Duration(0)).Return(fmt.Errorf("store: %w", repository.ErrKeyExisted)).Once()
				storage.EXPECT().StoreURL(ctx, "free000", url, time.Duration(0)).Return(nil).Once()
			},
			expectedCode: "free000",
		},
		{
			name: "gives up after max retries",
			setup: func(ctx context.Context, gen *mocks.GenPass, storage *repomocks.URLStorage) {
				gen.EXPECT().GeneratePassword(DefaultCodeLength).Return("taken00", nil).Times(MaxGenerateRetries)
				storage.EXPECT().StoreURL(ctx, "taken00", url, time.Duration(0)).Return(repository.ErrKeyExisted).Times(MaxGenerateRetries)
			},
			expectedErr: ErrCodeExhausted,
		},
		{
			name: "code generator error",
			setup: func(_ context.Context, gen *mocks.GenPass, _ *repomocks.URLStorage) {
				gen.EXPECT().GeneratePassword(DefaultCodeLength).Return("", errors.New("rand failure")).Once()
			},
			expectedErr: errors.New("rand failure"),
		},
		{
			name: "storage error",
			setup: func(ctx context.Context, gen *mocks.GenPass, storage *repomocks.URLStorage) {
				gen.EXPECT().GeneratePassword(DefaultCodeLength).Return("abc1234", nil).Once()
				storage.EXPECT().StoreURL(ctx, "abc1234", url, time.Duration(0)).Return(errors.New("redis down")).Once()
			},
			expectedErr: errors.New("redis down"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gen := mocks.NewGenPass(t)
			storage := repomocks.NewURLStorage(t)
			ctx := newRequestContext()
			tc.setup(ctx, gen, storage)

			code, err := NewShortenUrl(storage, gen).CreateShortenLink(ctx, url, tc.exp)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedCode, code)
		})
	}
}
