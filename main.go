// Package traefik_enforce_header_case_plugin A Traefik middleware plugin which enforces the case of specified request headers.
package traefik_enforce_header_case_plugin

import (
	"context"
	"net/http"
	"net/textproto"
)

// Config plugin configuration.
type Config struct {
	Headers []string `json:"headers,omitempty"`
}

// CreateConfig creates an empty plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: []string{},
	}
}

// TraefikEnforceHeaderCaseMiddleware plugin.
type TraefikEnforceHeaderCaseMiddleware struct {
	name    string
	next    http.Handler
	headers []string
}

// New creates a new TraefikEnforceHeaderCaseMiddleware plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &TraefikEnforceHeaderCaseMiddleware{
		name:    name,
		next:    next,
		headers: config.Headers,
	}, nil
}

func (ehc *TraefikEnforceHeaderCaseMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, key := range ehc.headers {
		value, ok := req.Header[textproto.CanonicalMIMEHeaderKey(key)]

		if ok {
			req.Header.Del(key)
			req.Header[key] = value
		}
	}

	ehc.next.ServeHTTP(rw, req)
}
