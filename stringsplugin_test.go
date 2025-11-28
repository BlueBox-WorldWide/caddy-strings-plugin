package stringsplugin

import (
	"context" // Required for context.WithValue
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// nextHandler satisfies the caddyhttp.Handler interface and is used to verify
// that the Strings plugin correctly modified the request replacer.
type nextHandler struct {
	t *testing.T
}

// ServeHTTP attempts to use the custom placeholders added by the Strings module
// and writes the result to the response body for verification.
func (h nextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	// 1. Retrieve the replacer from the context
	replVal := r.Context().Value(caddy.ReplacerCtxKey)
	if replVal == nil {
		h.t.Fatal("FATAL: Replacer not found in context (plugin didn't pass through or context was lost)")
	}
	repl := replVal.(*caddy.Replacer)

	// 2. Test the .lower mapping
	// The original value is "My-Test-String"
	lowerVal := repl.ReplaceAll("{test_var.lower}", "")
	if lowerVal != "my-test-string" {
		h.t.Errorf("FAILURE: .lower mapping failed. Expected 'my-test-string', got '%s'", lowerVal)
	}

	// 3. Test the .upper mapping
	upperVal := repl.ReplaceAll("{test_var.upper}", "")
	if upperVal != "MY-TEST-STRING" {
		h.t.Errorf("FAILURE: .upper mapping failed. Expected 'MY-TEST-STRING', got '%s'", upperVal)
	}
    
    // 4. Test non-existent variable (should return an empty string in this ReplaceAll context)
    // The previous expectation was incorrect for the common behavior of ReplaceAll("{var}", "")
    missingVal := repl.ReplaceAll("{missing_var.lower}", "")
    if missingVal != "" {
        h.t.Errorf("FAILURE: Missing variable lookup failed. Expected '', got '%s'", missingVal)
    }

	// If all checks pass, write a success marker.
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Replacer and pass-through successful"))
	return nil
}

func TestServeHTTP_AppliesReplacer(t *testing.T) {
	m := &Strings{}

	// Create a base replacer with the value we want to manipulate
	repl := caddy.NewReplacer()
	repl.Map(func(key string) (any, bool) {
		if key == "test_var" {
			return "My-Test-String", true
		}
		return nil, false
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.test/", nil)

	// Inject the mock replacer into the request context, simulating Caddy's framework setup
	ctx := context.WithValue(req.Context(), caddy.ReplacerCtxKey, repl)
	req = req.WithContext(ctx)

	// Provision is called for completeness, though it does nothing here
	if err := m.Provision(caddy.Context{}); err != nil {
		t.Fatalf("Provision returned error: %v", err)
	}

	// Call the ServeHTTP handler
	err := m.ServeHTTP(rr, req, nextHandler{t: t})
	if err != nil {
		t.Fatalf("ServeHTTP returned error: %v", err)
	}

	// Check if the response was successful (indicating the next handler was reached)
	expectedBody := "Replacer and pass-through successful"
	if !strings.Contains(rr.Body.String(), expectedBody) {
		t.Fatalf("Expected response body to contain %q, got: %q", expectedBody, rr.Body.String())
	}
}

func TestProvision_DoesNotError(t *testing.T) {
	m := &Strings{}
	if err := m.Provision(caddy.Context{}); err != nil {
		t.Fatalf("Provision returned error: %v", err)
	}
}

// Ensure Strings satisfies the handler-related interfaces used by Caddy
var (
	_ caddy.Module                = (*Strings)(nil)
	_ caddy.Provisioner           = (*Strings)(nil)
	_ caddyhttp.MiddlewareHandler = (*Strings)(nil)
)
