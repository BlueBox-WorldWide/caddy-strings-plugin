package strings_plugin

import (
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(CaddyStrings{})
	httpcaddyfile.RegisterHandlerDirective("strings", parseCaddyfile)
	httpcaddyfile.RegisterDirectiveOrder("strings", "before", "redir")
}

// Interface guards
var (
	_ caddy.Module                = (*CaddyStrings)(nil)
	_ caddy.Provisioner           = (*CaddyStrings)(nil)
	_ caddyhttp.MiddlewareHandler = (*CaddyStrings)(nil)
	_ caddyfile.Unmarshaler       = (*CaddyStrings)(nil)
)

type CaddyStrings struct{}

// CaddyModule returns module info
func (CaddyStrings) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.strings",
		New: func() caddy.Module { return new(CaddyStrings) },
	}
}

// Provision registers a global mapping for placeholders
func (m *CaddyStrings) Provision(ctx caddy.Context) error {
    return nil
}

func (m *CaddyStrings) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	
	mapStringCases(repl)

	return next.ServeHTTP(w, r)
}

//function that takes repl and maps .lower and .upper
func mapStringCases(repl *caddy.Replacer) {
	repl.Map(func(key string) (any, bool) {
		base := key
		lower := false
		upper := false

		if strings.HasSuffix(key, ".lower") {
			base = strings.TrimSuffix(key, ".lower")
			lower = true
		} else if strings.HasSuffix(key, ".upper") {
			base = strings.TrimSuffix(key, ".upper")
			upper = true
		}

		if lower || upper {
			val, found := repl.GetString(base)
			if !found {
				return nil, false
			}
			if lower {
				val = strings.ToLower(val)
			} else if upper {
				val = strings.ToUpper(val)
			}
			return val, true
		}

		return nil, false
	})
}

// UnmarshalCaddyfile reads static options (none for now)
func (m *CaddyStrings) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

// Caddyfile parser
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m CaddyStrings
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return &m, err
}
