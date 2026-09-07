package gateway

import (
"net/http"
"strings"

"./decoy"
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
// Agent request → handle by C2 (langsung forward ke handler yg ada)
return
case VisitorOperator:
// Operator request → handle by dashboard (nanti)
return
case VisitorScanner, VisitorUnknown:
// Scanner atau visitor tanpa header → decoy page
decoy.ServeDecoy(w, r)
}
}
