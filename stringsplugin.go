package stringsplugin

import (
	"fmt"
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
	_ caddy.Provisioner           = (*Strings)(nil)
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

// Provision registers a global mapping for placeholders
func (m *Strings) Provision(ctx caddy.Context) error {
    // global mapping for .upper/.lower
    repl := caddy.NewReplacer()
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

        v, ok := repl.Get(base)
        if !ok || v == nil {
            return nil, false
        }
        val := fmt.Sprintf("%v", v)
        if lower {
            val = strings.ToLower(val)
        } else if upper {
            val = strings.ToUpper(val)
        }
        return val, true
    })

    return nil
}

// UnmarshalCaddyfile reads static options (none for now)
func (m *Strings) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}

// ServeHTTP registers the mapping for per-request placeholders
func (m *Strings) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	repl, ok := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	if ok {
		repl.Map(stringsMappingFunc)
	}
	return next.ServeHTTP(w, r)
}

// Mapping function for .lower and .upper
func stringsMappingFunc(key string) (any, bool) {
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

	v, ok := caddy.NewReplacer().Get(base) // get global or per-request value
	if !ok || v == nil {
		return nil, false
	}

	val := fmt.Sprintf("%v", v)
	if lower {
		val = strings.ToLower(val)
	} else if upper {
		val = strings.ToUpper(val)
	}
	return val, true
}

// Caddyfile parser
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Strings
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return &m, err
}
