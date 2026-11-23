package service

import (
	"context"
	"time"

	"github.com/phihdn/golink/internal/model"
	"github.com/phihdn/golink/internal/repository"
)

// pingTimeout caps how long the health check waits for Redis, so an unreachable host fails fast.
const pingTimeout = time.Second

type Health interface {
	// Check returns the service info, or an error when Redis does not answer a ping.
	Check(ctx context.Context) (model.Health, error)
}

type health struct {
	serviceName string
	instanceID  string
	r           repository.HealthCheck
}

func NewHealth(serviceName, instanceID string, r repository.HealthCheck) Health {
	return &health{serviceName: serviceName, instanceID: instanceID, r: r}
}

func (h *health) Check(ctx context.Context) (model.Health, error) {
	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := h.r.Ping(ctx); err != nil {
		return model.Health{}, err
	}

	return model.Health{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}, nil
}
