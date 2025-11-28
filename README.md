# Caddy Strings Plugin

A small Caddy v2 module that extends the internal Caddy replacer, allowing you to apply simple string manipulations (like `.upper` and `.lower`) to any existing placeholder dynamically.

## Features

- `.lower` — convert a placeholder value to lowercase.
- `.upper` — convert a placeholder value to uppercase.

## Installation

Build Caddy with this module using xcaddy:

```bash
xcaddy build \
  --with github.com/BlueBox-WorldWide/caddy-strings-plugin
```

## Caddyfile — ordering

This third-party middleware must be ordered so it runs before the directive that consumes the modified placeholders (for example `respond`, `rewrite`, or `reverse_proxy`).

### Example — simple respond

Global options must order the middleware:

```caddyfile
{
    order strings before respond
}

:8080 {
    strings

    # Request: curl "http://localhost:8080/?name=caddy"
    # Output: HELLO CADDY!
    respond "HELLO {query.name.upper}!"
}
```

### Example — reverse proxy (normalize path)

```caddyfile
{
    order strings before reverse_proxy
}

:8080 {
    strings

    # Always send the path to the upstream in lowercase
    rewrite * {path.lower}
    reverse_proxy localhost:9000
}
```

## JSON configuration example

```json
{
  "apps": {
    "http": {
      "servers": {
        "srv0": {
          "listen": [":8080"],
          "routes": [
            {
              "handle": [
                { "handler": "strings" },
                {
                  "handler": "static_response",
                  "body": "Method: {http.request.method.lower}"
                }
              ]
            }
          ]
        }
      }
    }
  }
}
```

## Development

1. Clone the repository.
2. Run Caddy with the plugin loaded:

```bash
xcaddy run
```

## Notes

- Ensure the `order strings before ...` directive is present in global options when using handlers that rely on modified placeholders.
- This README assumes the standard Caddy replacer placeholders (e.g., `{query.name}`, `{path}`).
