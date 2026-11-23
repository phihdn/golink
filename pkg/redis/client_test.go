package redis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_ConnectsWithEnvSettings(t *testing.T) {
	mr := miniredis.RunT(t)
	t.Setenv("REDIS_ADDR", mr.Addr())
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "3")

	client, err := NewClient("")
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	assert.Equal(t, mr.Addr(), client.Options().Addr)
	assert.Equal(t, 3, client.Options().DB)
	require.NoError(t, client.Ping(context.Background()).Err())
}

func TestNewClient_InvalidConfig(t *testing.T) {
	t.Setenv("REDIS_DB", "not-a-number")

	client, err := NewClient("")

	assert.Error(t, err)
	assert.Nil(t, client)
}
