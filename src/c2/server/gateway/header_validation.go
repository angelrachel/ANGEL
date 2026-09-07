package gateway

import (
"os"
"strings"
)

const (
HeaderOperatorKey = "X-Operator-Key"
HeaderBeaconToken = "X-Beacon-Token"
)

func GetOperatorKey() string {
return os.Getenv("ANGEL_OPERATOR_KEY")
}

func ValidateOperatorKey(key string) bool {
expected := GetOperatorKey()
if expected == "" {
return true
}
return strings.EqualFold(key, expected)
}

func IsBeaconRequest(headers map[string][]string) bool {
for key, values := range headers {
if strings.EqualFold(key, HeaderBeaconToken) {
if len(values) > 0 && values[0] != "" {
return true
}
}
}
return false
}

func IsOperatorRequest(headers map[string][]string) bool {
for key, values := range headers {
if strings.EqualFold(key, HeaderOperatorKey) {
if len(values) > 0 && values[0] != "" {
return ValidateOperatorKey(values[0])
}
}
}
return false
}
