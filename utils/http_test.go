package utils

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeReverseProxyForwardedHeaders(t *testing.T) {
	var gotProto, gotHost string
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotProto = r.Header.Get("X-Forwarded-Proto")
		gotHost = r.Header.Get("X-Forwarded-Host")
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	rp, err := MakeReverseProxy(backend.URL, false, 0)
	if err != nil {
		t.Fatalf("MakeReverseProxy failed: %v", err)
	}

	t.Run("sets https when TLS is present", func(t *testing.T) {
		gotProto, gotHost = "", ""
		req := httptest.NewRequest(http.MethodGet, "https://example.com/test", nil)
		req.TLS = &tls.ConnectionState{}
		rec := httptest.NewRecorder()
		rp.ServeHTTP(rec, req)
		if gotProto != "https" {
			t.Errorf("expected X-Forwarded-Proto=https, got %q", gotProto)
		}
		if gotHost != "example.com" {
			t.Errorf("expected X-Forwarded-Host=example.com, got %q", gotHost)
		}
	})

	t.Run("preserves existing X-Forwarded-Proto", func(t *testing.T) {
		gotProto, gotHost = "", ""
		req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
		req.Header.Set("X-Forwarded-Proto", "https")
		rec := httptest.NewRecorder()
		rp.ServeHTTP(rec, req)
		if gotProto != "https" {
			t.Errorf("expected X-Forwarded-Proto=https, got %q", gotProto)
		}
		if gotHost != "example.com" {
			t.Errorf("expected X-Forwarded-Host=example.com, got %q", gotHost)
		}
	})

	t.Run("sets http when plain HTTP", func(t *testing.T) {
		gotProto, gotHost = "", ""
		req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
		rec := httptest.NewRecorder()
		rp.ServeHTTP(rec, req)
		if gotProto != "http" {
			t.Errorf("expected X-Forwarded-Proto=http, got %q", gotProto)
		}
	})
}
