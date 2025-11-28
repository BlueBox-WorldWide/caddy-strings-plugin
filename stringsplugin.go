package stringsplugin

import (
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Strings{})
	httpcaddyfile.RegisterHandlerDirective("strings", parseCaddyfile)
	httpcaddyfile.RegisterDirectiveOrder("strings", "before", "redir")
}

// Interface guards
var (
	_ caddy.Module                = (*Strings)(nil)
	_ caddy.Provisioner           = (*Strings)(nil)
	_ caddyhttp.MiddlewareHandler = (*Strings)(nil)
)

type Strings struct {
	// Message may contain placeholders (global or request-specific).
	Message string `json:"message,omitempty"`
}

func (Strings) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.strings",
		New: func() caddy.Module { return new(Strings) },
	}
}

// Provision implements caddy.Provisioner.
func (m *Strings) Provision(ctx caddy.Context) error {
	return nil
}

// Validate implements caddy.Validator.
func (m *Strings) Validate() error {
	return nil
}

func (m *Strings) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// 1. Get the current replacer (has access to request-time placeholders)
	reqRepl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)

	// 2. Inject the mapping function so .lower/.upper also work on request-time placeholders
	reqRepl.Map(func(key string) (any, bool) {
		// Handle .lower
		if strings.HasSuffix(key, ".lower") {
			base := strings.TrimSuffix(key, ".lower")
			if v, ok := reqRepl.GetString(base); ok {
				return strings.ToLower(v), true
			}
		}

		// Handle .upper
		if strings.HasSuffix(key, ".upper") {
			base := strings.TrimSuffix(key, ".upper")
			if v, ok := reqRepl.GetString(base); ok {
				return strings.ToUpper(v), true
			}
		}

		return nil, false
	})

	// 3. Resolve any remaining placeholders in the configured message using the request replacer,
	//    then write it to the response.
	resolved := reqRepl.ReplaceAll(m.Message, "")
	if len(resolved) > 0 {
		if _, err := w.Write([]byte(resolved)); err != nil {
			return err
		}
	}

	// 4. Continue the chain
	return next.ServeHTTP(w, r)
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	h.Next() // consume directive
	return &Strings{}, nil
}
