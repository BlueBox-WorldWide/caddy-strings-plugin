package stringsplugin

import (
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
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
	_ caddyhttp.MiddlewareHandler = (*Strings)(nil)
	_ caddyfile.Unmarshaler       = (*Strings)(nil)
)

type Strings struct{}

// CaddyModule returns module info
func (Strings) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.strings",
		New: func() caddy.Module { return new(Strings) },
	}
}

func (m *Strings) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	repl, ok := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	if ok {
		repl.Map(func(key string) (any, bool) {
			lower := false
			upper := false
			base := key

			if strings.HasSuffix(key, ".lower") {
				base = strings.TrimSuffix(key, ".lower")
				lower = true
			} else if strings.HasSuffix(key, ".upper") {
				base = strings.TrimSuffix(key, ".upper")
				upper = true
			}

			val, _ := repl.ReplaceFunc(base, func(variable string, v any) (any, error) {
				str, ok := v.(string)
				if !ok {
					return v, nil // leave non-string values unchanged
				}
			
				if lower {
					return strings.ToLower(str), nil
				} else if upper {
					return strings.ToUpper(str), nil
				}
				return str, nil
			})			

			if val == "" {
				return "", false
			}
			return val, true
		})
	}

	return next.ServeHTTP(w, r)
}


// UnmarshalCaddyfile reads static options (none for now)
func (m *Strings) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

// Caddyfile parser
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Strings
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return &m, err
}
