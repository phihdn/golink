package redis

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// unsetEnv removes the variables for this test and restores them afterwards.
func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, k := range keys {
		t.Setenv(k, "") // registers the restore on cleanup
		require.NoError(t, os.Unsetenv(k))
	}
}

func TestNewConfig_Defaults(t *testing.T) {
	unsetEnv(t, "REDIS_ADDR", "REDIS_PASSWORD", "REDIS_DB")

	cfg, err := newConfig("")
	require.NoError(t, err)

	assert.Equal(t, &config{Address: "localhost:6379", Password: "", DB: 0}, cfg)
}

func TestNewConfig_FromEnv(t *testing.T) {
	t.Setenv("REDIS_ADDR", "redis:6380")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("REDIS_DB", "2")

	cfg, err := newConfig("")
	require.NoError(t, err)

	assert.Equal(t, &config{Address: "redis:6380", Password: "secret", DB: 2}, cfg)
}

func TestNewConfig_PrefixWinsOverPlainName(t *testing.T) {
	unsetEnv(t, "CACHE_REDIS_PASSWORD", "CACHE_REDIS_DB", "REDIS_PASSWORD", "REDIS_DB")
	t.Setenv("REDIS_ADDR", "plain:6379")
	t.Setenv("CACHE_REDIS_ADDR", "prefixed:6379")

	cfg, err := newConfig("CACHE")
	require.NoError(t, err)

	assert.Equal(t, "prefixed:6379", cfg.Address)
}

func TestNewConfig_PrefixFallsBackToPlainName(t *testing.T) {
	unsetEnv(t, "CACHE_REDIS_ADDR", "CACHE_REDIS_PASSWORD", "CACHE_REDIS_DB", "REDIS_PASSWORD")
	t.Setenv("REDIS_ADDR", "plain:6379")
	t.Setenv("REDIS_DB", "4")

	cfg, err := newConfig("CACHE")
	require.NoError(t, err)

	assert.Equal(t, &config{Address: "plain:6379", Password: "", DB: 4}, cfg)
}

func TestNewConfig_InvalidDB(t *testing.T) {
	t.Setenv("REDIS_DB", "not-a-number")

	cfg, err := newConfig("")

	assert.Error(t, err)
	assert.Nil(t, cfg)
}
