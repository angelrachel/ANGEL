package api

import (
"net/http"
"time"
)

type Middleware func(http.Handler) http.Handler

func AuthMiddleware(token string) Middleware {
return func(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.Header.Get("Authorization") != "Bearer "+token {
http.Error(w, "Unauthorized", http.StatusUnauthorized)
return
}
next.ServeHTTP(w, r)
})
}
}

func LoggingMiddleware() Middleware {
return func(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
start := time.Now()
next.ServeHTTP(w, r)
_ = start
})
}
}

func Wrap(handler http.Handler, middlewares ...Middleware) http.Handler {
var final http.Handler = handler
for i := len(middlewares) - 1; i >= 0; i-- {
final = middlewares[i](final)
}
return final
}
