package stringsplugin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
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

type Strings struct {
	// logger provides structured logging for the plugin's internal operations.
	logger *zap.Logger
}

// CaddyModule returns module info
func (Strings) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.strings",
		New: func() caddy.Module { return new(Strings) },
	}
}

func (m *Strings) Provision(ctx caddy.Context) error {
	m.logger = ctx.Logger() // m.logger is a *zap.Logger

	return nil
}

func (m *Strings) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// log request start
	m.logger.Debug("strings handler ServeHTTP start",
		zap.String("method", r.Method),
		zap.String("uri", r.URL.String()),
	)

	repl, ok := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	if ok {
		if m.logger != nil {
			m.logger.Debug("replacer found in request context")
		}

		repl.Map(func(key string) (any, bool) {
			if m.logger != nil {
				m.logger.Debug("processing replacement key", zap.String("key", key))
			}

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

			if m.logger != nil {
				m.logger.Debug("derived base/modifier", zap.String("base", base), zap.Bool("lower", lower), zap.Bool("upper", upper))
			}

			val, _ := repl.ReplaceFunc(base, func(variable string, v any) (any, error) {
				if m.logger != nil {
					m.logger.Debug("ReplaceFunc called", zap.String("variable_requested", variable))
				}

				str, ok := v.(string)
				if !ok {
					if m.logger != nil {
						m.logger.Debug("value is not a string; leaving unchanged", zap.String("value_type", fmt.Sprintf("%T", v)), zap.Any("raw_value", v))
					}
					return v, nil // leave non-string values unchanged
				}

				if lower {
					if m.logger != nil {
						m.logger.Debug("applying lower", zap.String("original", str))
					}
					return strings.ToLower(str), nil
				} else if upper {
					if m.logger != nil {
						m.logger.Debug("applying upper", zap.String("original", str))
					}
					return strings.ToUpper(str), nil
				}

				if m.logger != nil {
					m.logger.Debug("returning original string", zap.String("value", str))
				}
				return str, nil
			})

			if val == "" {
				if m.logger != nil {
					m.logger.Debug("replacement produced empty string", zap.String("key", key), zap.Any("value", val))
				}
				return "", false
			}
			if m.logger != nil {
				m.logger.Debug("replacement produced value", zap.String("key", key), zap.Any("value", val))
			}
			return val, true
		})
	} else {
		if m.logger != nil {
			m.logger.Debug("no replacer present in request context")
		}
	}

	if m.logger != nil {
		m.logger.Debug("calling next handler")
	}
	err := next.ServeHTTP(w, r)
	if err != nil {
		if m.logger != nil {
			m.logger.Error("next handler returned error", zap.Error(err))
		}
	} else {
		if m.logger != nil {
			m.logger.Debug("next handler completed without error")
		}
	}
	return err
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
