package traefik_enforce_header_case_plugin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clugg/traefik-enforce-header-case-plugin"
)

func preparePlugin(t *testing.T) (context.Context, *traefik_enforce_header_case_plugin.Config, http.Handler) {
	t.Helper()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	cfg := traefik_enforce_header_case_plugin.CreateConfig()
	cfg.Headers = append(cfg.Headers, "x-tEsT-hEAder")

	handler, err := traefik_enforce_header_case_plugin.New(ctx, next, cfg, "traefik-enforce-header-case-plugin")
	if err != nil {
		t.Fatal(err)
	}

	return ctx, cfg, handler
}

func prepareRequest(ctx context.Context, t *testing.T) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	return recorder, req
}

func TestEnforceHeaderCaseWhenHeaderExists(t *testing.T) {
	ctx, cfg, handler := preparePlugin(t)
	recorder, req := prepareRequest(ctx, t)

	req.Header.Set("x-test-header", "123")
	handler.ServeHTTP(recorder, req)

	canonicalHeaderValue := req.Header.Get(cfg.Headers[0])
	if canonicalHeaderValue != "" {
		t.Errorf("unexpected value for canonicalised header: %s", canonicalHeaderValue)
	}

	caseEnforcedHeaderValue, ok := req.Header[cfg.Headers[0]]
	if !ok || len(caseEnforcedHeaderValue) != 1 || caseEnforcedHeaderValue[0] != "123" {
		t.Errorf("unexpected value for case-enforced header: %q", caseEnforcedHeaderValue)
	}
}

func TestEnforceHeaderCaseWhenHeaderDoesNotExist(t *testing.T) {
	ctx, cfg, handler := preparePlugin(t)
	recorder, req := prepareRequest(ctx, t)

	handler.ServeHTTP(recorder, req)

	canonicalHeaderValue := req.Header.Get(cfg.Headers[0])
	if canonicalHeaderValue != "" {
		t.Errorf("unexpected value for canonicalised header: %s", canonicalHeaderValue)
	}

	caseEnforcedHeaderValue, ok := req.Header[cfg.Headers[0]]
	if ok || len(caseEnforcedHeaderValue) != 0 {
		t.Errorf("unexpected value for case-enforced header: %q", caseEnforcedHeaderValue)
	}
}

func TestEnforceHeaderCaseWhenUnconfiguredHeaderExists(t *testing.T) {
	ctx, _, handler := preparePlugin(t)
	recorder, req := prepareRequest(ctx, t)

	req.Header.Set("x-test-header-copy", "123")
	handler.ServeHTTP(recorder, req)

	canonicalHeaderValue := req.Header.Get("x-test-header-copy")
	if canonicalHeaderValue != "123" {
		t.Errorf("unexpected value for unconfigured header: %s", canonicalHeaderValue)
	}
}
