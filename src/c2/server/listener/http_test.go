package listener

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHTTPServerConfiguresTimeoutsAndHandler(t *testing.T) {
	server := NewHTTPServer(8081)
	if server.Port != 8081 || server.Server.Addr != ":8081" {
		t.Fatalf("unexpected server configuration: %+v", server)
	}
	if server.Server.ReadHeaderTimeout <= 0 || server.Server.ReadTimeout <= 0 || server.Server.WriteTimeout <= 0 || server.Server.IdleTimeout <= 0 {
		t.Fatal("expected all HTTP timeouts to be configured")
	}
	server.RegisterHandler("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	server.Server.Handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d", recorder.Code)
	}
}

func TestListenAndServeReturnsAfterCancellation(t *testing.T) {
	server := NewHTTPServer(0)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	started := make(chan error, 1)
	go func() { started <- server.ListenAndServe(ctx) }()
	select {
	case err := <-started:
		if err == nil {
			t.Fatal("expected bind error for invalid port address")
		}
	case <-time.After(time.Second):
		t.Fatal("listener did not return")
	}
}
