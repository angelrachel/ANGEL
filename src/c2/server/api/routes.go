package api

import (
	"net/http"
	"sync"
	"time"
)

type Route struct {
	Pattern string
	Method  string
	Handler http.HandlerFunc
}

type Router struct {
	mu     sync.RWMutex
	Routes []Route
}

func NewRouter() *Router {
	return &Router{Routes: []Route{}}
}

func (r *Router) Register(route Route) {
	if route.Pattern == "" || route.Method == "" || route.Handler == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Routes = append(r.Routes, route)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	routes := append([]Route(nil), r.Routes...)
	r.mu.RUnlock()
	for _, route := range routes {
		if req.URL.Path == route.Pattern && req.Method == route.Method {
			route.Handler(w, req)
			return
		}
	}
	http.NotFound(w, req)
}

func (r *Router) List() []Route {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Route(nil), r.Routes...)
}

func (r *Router) HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().UTC().Format(time.RFC3339) + `"}`))
}
