package config

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("SERVICE_NAME", "")
	t.Setenv("INSTANCE_ID", "")

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "golink", cfg.ServiceName)
	_, parseErr := uuid.Parse(cfg.InstanceID)
	assert.NoError(t, parseErr, "instance id should be a generated UUID")
}

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("SERVICE_NAME", "my_service")
	t.Setenv("INSTANCE_ID", "instance-1")

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, Config{Port: "9090", ServiceName: "my_service", InstanceID: "instance-1"}, cfg)
}
