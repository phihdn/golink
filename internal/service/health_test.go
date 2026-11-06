package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/phihdn/golink/internal/model"
)

func TestHealth_Check(t *testing.T) {
	svc := NewHealth("golink", "abc-123")

	assert.Equal(t, model.Health{
		Message:     "OK",
		ServiceName: "golink",
		InstanceID:  "abc-123",
	}, svc.Check())
}
