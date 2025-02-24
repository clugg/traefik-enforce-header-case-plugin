# Traefik Enforce Header Case Plugin

![Build Status](https://github.com/clugg/traefik-enforce-header-case-plugin/workflows/Main/badge.svg) [![Latest Release](https://img.shields.io/github/v/release/clugg/traefik-enforce-header-case-plugin?include_prereleases&sort=semver)](https://github.com/clugg/traefik-enforce-header-case-plugin/releases)

A [Traefik](https://traefik.io/) middleware plugin which enforces the case of specified request headers.

This overrides Go (and, by extension, Traefik)'s default behaviour of [canonicalising header keys](https://pkg.go.dev/net/http#Header), which can be useful when working with HTTP servers/applications that are case-sensitive to certain headers.

## Configuration

### Static

```yaml
experimental:
  plugins:
    enforceHeaderCase:
      moduleName: github.com/clugg/traefik-enforce-header-case-plugin
      version: v0.1.0
```

### Dynamic

In the following example, a router has been set up with middleware set up to ensure that the `x-tEsT-hEAder` header is always forwaded to the service exactly as it appears in the configuration. This means that if a request containing this header is made to the router with any variation of casing (e.g. `x-test-header: 123` or `X-TEST-HEADER: 123`), the request forwarded to the service will contain `x-tEsT-hEAder: 123`.

```yaml
http:
  middlewares:
    enforce-header-case:
      plugin:
        enforceHeaderCase:
          headers:
            - x-tEsT-hEAder

  routers:
    my-router:
      rule: Host(`localhost`)
      service: my-service
      middlewares:
        - enforce-header-case@file

  services:
    my-service:
      loadBalancer:
        servers:
          - url: 'http://127.0.0.1'
```
