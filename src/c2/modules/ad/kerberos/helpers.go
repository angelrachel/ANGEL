package kerberos

const base64Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

func base64Encode(data []byte) string {
result := ""
bits := 0
val := 0
for _, b := range data {
val = (val << 8) | int(b)
bits += 8
for bits >= 6 {
bits -= 6
result += string(base64Alphabet[(val>>bits)&0x3F])
}
}
if bits > 0 {
result += string(base64Alphabet[(val<<(6-bits))&0x3F])
}
for len(result)%4 != 0 {
result += "="
}
return result
}
