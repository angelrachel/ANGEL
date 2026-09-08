package api

import (
"net/http"
"time"
)

type Route struct {
Pattern string
Method  string
Handler http.HandlerFunc
}

type Router struct {
Routes []Route
}

func NewRouter() *Router {
return &Router{Routes: []Route{}}
}

func (r *Router) Register(route Route) {
r.Routes = append(r.Routes, route)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
for _, route := range r.Routes {
if req.URL.Path == route.Pattern && req.Method == route.Method {
route.Handler(w, req)
return
}
}
http.NotFound(w, req)
}

func (r *Router) List() []Route {
return r.Routes
}

func (r *Router) HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().UTC().Format(time.RFC3339) + `"}`))
}
