package stringsplugin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caddyserver/caddy/v2"
)

type nextHandler struct{}

func (nextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("next called"))
    return nil
}

func TestServeHTTP_PassesToNext(t *testing.T) {
    m := &Strings{}

    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "http://example.test/", nil)

    err := m.ServeHTTP(rr, req, nextHandler{})
    if err != nil {
        t.Fatalf("ServeHTTP returned error: %v", err)
    }

    if rr.Body.String() != "next called" {
        t.Fatalf("unexpected response body: %q", rr.Body.String())
    }
}

func TestProvision_DoesNotError(t *testing.T) {
    m := &Strings{}
    if err := m.Provision(caddy.Context{}); err != nil {
        t.Fatalf("Provision returned error: %v", err)
    }
}

// Ensure Strings satisfies the handler-related interfaces used by Caddy
var _ caddy.Module = (*Strings)(nil)
// Note: Do not assert `caddyhttp.Middleware` here to remain compatible
// with different Caddy versions which may not expose that exact type.
