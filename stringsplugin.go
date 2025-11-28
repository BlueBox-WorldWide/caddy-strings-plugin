package stringsplugin

// Package stringsplugin provides a Caddy module that implements
// placeholder filters for string manipulation, specifically
// converting strings to lowercase or uppercase. It registers
// these filters so they can be used in Caddy's configuration.

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

// Strings provides lowercase/uppercase placeholder filters.
type Strings struct{}

func (Strings) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.strings",
		New: func() caddy.Module { return new(Strings) },
	}
}

func (m *Strings) Provision(ctx caddy.Context) error {
	caddyhttp.SetVarReplacementFunc("lower", func(_ *caddyhttp.RequestContext, in string) (string, error) {
		return strings.ToLower(in), nil
	})

	caddyhttp.SetVarReplacementFunc("upper", func(_ *caddyhttp.RequestContext, in string) (string, error) {
		return strings.ToUpper(in), nil
	})

	return nil
}

func (m *Strings) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	return next.ServeHTTP(w, r)
}

func (m *Strings) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	return nil
}
