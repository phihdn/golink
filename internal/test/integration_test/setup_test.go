package integration_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/phihdn/golink/internal/api"
	"github.com/phihdn/golink/internal/config"
)

// newServer starts the full API against an in-memory Redis.
func newServer(t *testing.T, cfg config.Config) (api.Engine, *miniredis.Miniredis) {
	t.Helper()

	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = client.Close() })

	return api.NewEngine(&cfg, client), mr
}
