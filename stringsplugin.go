package stringsplugin

import (
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
    caddy.RegisterModule(Strings{})
}

// Interface guards
var (
    _ caddy.Module                = (*Strings)(nil)
    _ caddy.Provisioner           = (*Strings)(nil)
    _ caddyhttp.MiddlewareHandler = (*Strings)(nil)
    _ caddyfile.Unmarshaler       = (*Strings)(nil)
)

type Strings struct{}

func (Strings) CaddyModule() caddy.ModuleInfo {
    return caddy.ModuleInfo{
        ID:  "http.handlers.strings",
        New: func() caddy.Module { return new(Strings) },
    }
}

func (m *Strings) Provision(ctx caddy.Context) error {
    return nil
}

func (m *Strings) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
    return nil
}

func (m *Strings) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
    // 1. Get the current replacer
    repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)

    // 2. Inject the mapping function
    // This allows {upstream_host.lower} or {header.User-Agent.upper}
    repl.Map(func(key string) (any, bool) {
        // Handle .lower
        if strings.HasSuffix(key, ".lower") {
            base := strings.TrimSuffix(key, ".lower")
            
            // We must retrieve the value of the base key from the replacer
            if v, ok := repl.GetString(base); ok {
                return strings.ToLower(v), true
            }
        }
        
        // Handle .upper
        if strings.HasSuffix(key, ".upper") {
            base := strings.TrimSuffix(key, ".upper")
            
            if v, ok := repl.GetString(base); ok {
                return strings.ToUpper(v), true
            }
        }
        
        return nil, false
    })

    // 3. Continue the chain
    return next.ServeHTTP(w, r)
}
