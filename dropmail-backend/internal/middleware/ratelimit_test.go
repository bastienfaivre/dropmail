package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimiter_AllowsRequestsWithinLimit(t *testing.T) {
	rl := NewRateLimiter(60, 5) // 60 per minute, burst of 5

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First 5 requests should succeed (burst)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("request %d: expected status %d, got %d", i+1, http.StatusOK, rec.Code)
		}
	}
}

func TestRateLimiter_RejectsWhenLimitExceeded(t *testing.T) {
	rl := NewRateLimiter(60, 2) // 60 per minute, burst of 2

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Use up the burst
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "192.168.1.2:12345"
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}

	// Third request should be rate limited
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.2:12345"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("expected status %d, got %d", http.StatusTooManyRequests, rec.Code)
	}

	// Verify JSON error response
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["error"] != "Rate limit exceeded" {
		t.Errorf("unexpected error message: %s", resp["error"])
	}
}

func TestRateLimiter_DifferentIPsHaveSeparateLimits(t *testing.T) {
	rl := NewRateLimiter(60, 1) // 60 per minute, burst of 1

	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First IP uses up its limit
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.RemoteAddr = "192.168.1.3:12345"
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Errorf("first IP first request: expected %d, got %d", http.StatusOK, rec1.Code)
	}

	// First IP should be rate limited
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.RemoteAddr = "192.168.1.3:12345"
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusTooManyRequests {
		t.Errorf("first IP second request: expected %d, got %d", http.StatusTooManyRequests, rec2.Code)
	}

	// Second IP should still work
	req3 := httptest.NewRequest(http.MethodGet, "/", nil)
	req3.RemoteAddr = "192.168.1.4:12345"
	rec3 := httptest.NewRecorder()
	handler.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusOK {
		t.Errorf("second IP: expected %d, got %d", http.StatusOK, rec3.Code)
	}
}

func TestGetIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.100:54321"

	ip := getIP(req)
	if ip != "192.168.1.100" {
		t.Errorf("expected IP '192.168.1.100', got '%s'", ip)
	}
}

func TestGetIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.195")

	ip := getIP(req)
	if ip != "203.0.113.195" {
		t.Errorf("expected IP '203.0.113.195', got '%s'", ip)
	}
}

func TestGetIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	req.Header.Set("X-Real-IP", "198.51.100.178")

	ip := getIP(req)
	if ip != "198.51.100.178" {
		t.Errorf("expected IP '198.51.100.178', got '%s'", ip)
	}
}
