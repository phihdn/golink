package handler

import "context"

type ctxKey struct{}

// newRequestContext returns a context only this test owns. Mocks expect this
// exact value, so a layer that drops or replaces the caller's context fails the
// test. context.Background() alone would not catch that: every call returns the
// same value.
func newRequestContext() context.Context {
	return context.WithValue(context.Background(), ctxKey{}, "request-1")
}
