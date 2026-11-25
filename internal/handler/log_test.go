package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newLoggedContext returns a test context that carries a logger writing into
// the returned buffer, so a test can check what a handler logged.
func newLoggedContext(t *testing.T) (context.Context, *bytes.Buffer) {
	t.Helper()

	logs := &bytes.Buffer{}
	return zerolog.New(logs).WithContext(t.Context()), logs
}

// assertErrorLog checks that the handler logged exactly one error line with
// wantMsg and wantErr. An empty wantMsg means the handler must log nothing.
func assertErrorLog(t *testing.T, logs *bytes.Buffer, wantMsg, wantErr string) {
	t.Helper()

	if wantMsg == "" {
		assert.Empty(t, logs.String())
		return
	}

	var entry map[string]any
	require.NoError(t, json.Unmarshal(logs.Bytes(), &entry))
	assert.Equal(t, "error", entry["level"])
	assert.Equal(t, wantMsg, entry["message"])
	assert.Equal(t, wantErr, entry["error"])
}
