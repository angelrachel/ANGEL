package gateway

import (
"net/http"
"strings"

"ANGEL/src/c2/server/decoy"
)

type VisitorType int

const (
VisitorScanner VisitorType = iota
VisitorAgent
VisitorOperator
VisitorUnknown
)

func DetermineVisitorType(r *http.Request) VisitorType {
if IsBeaconRequest(r.Header) {
return VisitorAgent
}
if IsOperatorRequest(r.Header) {
return VisitorOperator
}
if strings.HasPrefix(r.URL.Path, "/login") {
return VisitorScanner
}
return VisitorScanner
}

func SelectRoute(w http.ResponseWriter, r *http.Request) {
vt := DetermineVisitorType(r)

switch vt {
case VisitorAgent:
return
case VisitorOperator:
return
case VisitorScanner, VisitorUnknown:
decoy.ServeDecoy(w, r)
}
}
