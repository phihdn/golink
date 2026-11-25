package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/phihdn/golink/internal/config"
)

func shorten(t *testing.T, srv http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

// decodeCode checks the success body and returns the generated code.
func decodeCode(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()

	require.Equal(t, http.StatusOK, rec.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "Shorten URL generated successfully!", resp["message"])
	assert.Regexp(t, `^[a-zA-Z0-9]{7}$`, resp["code"])
	return resp["code"]
}

func TestShorten_StoresURLWithClientExpiration(t *testing.T) {
	srv, mr := newServer(t, config.Config{})

	code := decodeCode(t, shorten(t, srv, `{"exp":604800,"url":"https://example.com/a/long/path"}`))

	stored, err := mr.Get(code)
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/a/long/path", stored)
	assert.Equal(t, 604800*time.Second, mr.TTL(code))

	mr.FastForward(604800 * time.Second)
	assert.False(t, mr.Exists(code), "link must be gone after exp")
}

func TestShorten_WithoutExpiration(t *testing.T) {
	srv, mr := newServer(t, config.Config{})

	code := decodeCode(t, shorten(t, srv, `{"url":"https://example.com"}`))

	assert.True(t, mr.Exists(code))
	assert.Equal(t, time.Duration(0), mr.TTL(code))
}

func TestShorten_InvalidInput(t *testing.T) {
	srv, _ := newServer(t, config.Config{})

	testCases := []struct {
		name string
		body string
	}{
		{name: "malformed json", body: `{"url":`},
		{name: "missing url", body: `{"exp":10}`},
		{name: "invalid url", body: `{"url":"not-a-url"}`},
		{name: "negative exp", body: `{"exp":-5,"url":"https://example.com"}`},
		{name: "exp too large", body: `{"exp":604801,"url":"https://example.com"}`},
		{name: "fractional exp", body: `{"exp":1.5,"url":"https://example.com"}`},
		{name: "string exp", body: `{"exp":"10","url":"https://example.com"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := shorten(t, srv, tc.body)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.JSONEq(t, `{"error":"Invalid Input"}`, rec.Body.String())
		})
	}
}

func TestShorten_RedisDown(t *testing.T) {
	srv, mr := newServer(t, config.Config{})
	mr.Close()

	rec := shorten(t, srv, `{"url":"https://example.com"}`)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"error":"Internal Server Error"}`, rec.Body.String())
}

func redirect(srv http.Handler, code string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/"+code, nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

func TestRedirect_ToStoredURL(t *testing.T) {
	srv, _ := newServer(t, config.Config{})

	code := decodeCode(t, shorten(t, srv, `{"url":"https://example.com/a/long/path"}`))
	rec := redirect(srv, code)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "https://example.com/a/long/path", rec.Header().Get("Location"))
}

func TestRedirect_UnknownCode(t *testing.T) {
	srv, _ := newServer(t, config.Config{})

	rec := redirect(srv, "missing")

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"error":"Not Found"}`, rec.Body.String())
}

func TestRedirect_ExpiredLink(t *testing.T) {
	srv, mr := newServer(t, config.Config{})

	code := decodeCode(t, shorten(t, srv, `{"exp":60,"url":"https://example.com"}`))
	mr.FastForward(60 * time.Second)
	rec := redirect(srv, code)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestRedirect_RedisDown(t *testing.T) {
	srv, mr := newServer(t, config.Config{})
	mr.Close()

	rec := redirect(srv, "abc1234")

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"error":"Internal Server Error"}`, rec.Body.String())
}
