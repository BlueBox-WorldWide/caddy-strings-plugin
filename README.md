# Caddy Strings Plugin

This repository provides a very small Caddy v2 module that registers two string placeholder replacement functions for use in Caddy configurations:

- `lower`: convert a string to lowercase
- `upper`: convert a string to uppercase

Usage

1. Build or include this module in a Caddy build that supports external modules (recommended to use `xcaddy`):

```
xcaddy build --with github.com/BlueBox-WorldWide/caddy-strings-plugin
```

2. After loading the module, you can use the replacers in your Caddyfile or other configs. Example (in a relevant context):

```
# Example usage depends on where you use replacers; this shows conceptually
{http.request.header.User-Agent|lower}
```

Notes

- This module only registers the replacer functions during `Provision` and otherwise is a no-op handler.
- The module ID is `http.handlers.strings`.
