package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// HealthCheck checks that the storage behind the app is reachable.
type HealthCheck interface {
	Ping(ctx context.Context) error
}

type healthCheck struct {
	c *redis.Client
}

func NewHealthCheck(c *redis.Client) HealthCheck {
	return &healthCheck{c: c}
}

func (h *healthCheck) Ping(ctx context.Context) error {
	return h.c.Ping(ctx).Err()
}
