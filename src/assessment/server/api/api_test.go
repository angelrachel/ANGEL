package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestAuthMiddlewareRejectsEmptyConfiguredToken(t *testing.T) {
	handler := AuthMiddleware("")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer ")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected empty configured token to fail closed, got %d", res.Code)
	}
}

func TestAuthMiddlewareAcceptsOnlyExactToken(t *testing.T) {
	handler := AuthMiddleware("operator-secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	for _, authorization := range []string{"", "Bearer wrong", "operator-secret", "Bearer operator-secret ", "Bearer operator-secret"} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", authorization)
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		if authorization == "Bearer operator-secret" && res.Code != http.StatusNoContent {
			t.Fatalf("expected exact token to pass, got %d", res.Code)
		}
		if authorization != "Bearer operator-secret" && res.Code != http.StatusUnauthorized {
			t.Fatalf("expected token %q to fail, got %d", authorization, res.Code)
		}
	}
}

func TestRouterRegisterListAndServeAreSafe(t *testing.T) {
	router := NewRouter()
	okHandler := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			router.Register(Route{Pattern: "/health", Method: http.MethodGet, Handler: okHandler})
			_ = router.List()
			res := httptest.NewRecorder()
			router.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/health", nil))
		}()
	}
	wg.Wait()

	if len(router.List()) == 0 {
		t.Fatal("expected registered routes")
	}
	if got := router.List(); strings.TrimSpace(got[0].Pattern) != "/health" {
		t.Fatalf("unexpected route: %#v", got[0])
	}
}

func TestHealthCheckSetsJSONContentType(t *testing.T) {
	router := NewRouter()
	res := httptest.NewRecorder()
	router.HealthCheckHandler(res, httptest.NewRequest(http.MethodGet, "/health", nil))
	if res.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	if !strings.HasPrefix(res.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("expected JSON content type, got %q", res.Header().Get("Content-Type"))
	}
}
