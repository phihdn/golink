package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenPass_GeneratePassword(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		expectedLen int
		expectedErr error
	}{
		{name: "zero length", expectedLen: 0, expectedErr: nil},
		{name: "short", expectedLen: 1, expectedErr: nil},
		{name: "normal", expectedLen: 16, expectedErr: nil},
		{name: "long", expectedLen: 256, expectedErr: nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := NewGenPass()
			password, err := svc.GeneratePassword(tc.expectedLen)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Len(t, password, tc.expectedLen)
			for _, r := range password {
				assert.True(t, strings.ContainsRune(charset, r), "unexpected character %q", r)
			}
		})
	}
}

func TestGenPass_GeneratePassword_IsRandom(t *testing.T) {
	t.Parallel()

	svc := NewGenPass()

	first, err := svc.GeneratePassword(32)
	require.NoError(t, err)
	second, err := svc.GeneratePassword(32)
	require.NoError(t, err)

	assert.NotEqual(t, first, second)
}
