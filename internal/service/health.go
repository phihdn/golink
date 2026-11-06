package service

import "github.com/phihdn/golink/internal/model"

type Health interface {
	Check() model.Health
}

type health struct {
	serviceName string
	instanceID  string
}

func NewHealth(serviceName, instanceID string) Health {
	return &health{serviceName: serviceName, instanceID: instanceID}
}

func (h *health) Check() model.Health {
	return model.Health{
		Message:     "OK",
		ServiceName: h.serviceName,
		InstanceID:  h.instanceID,
	}
}
