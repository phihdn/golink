package service

import (
	"context"

	"github.com/phihdn/golink/internal/model"
	"github.com/phihdn/golink/internal/repository"
)

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

// Check passes the caller's ctx to the repository unchanged, so the request
// context reaches Redis.
func (h *health) Check(ctx context.Context) (model.Health, error) {
	if err := h.r.Ping(ctx); err != nil {
		return model.Health{}, err
	}

	return model.Health{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}, nil
}
