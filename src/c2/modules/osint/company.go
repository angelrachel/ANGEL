package osint

import (
"encoding/json"
"net/http"
)

type CompanyOSINT struct{}

func NewCompanyOSINT() *CompanyOSINT {
return &CompanyOSINT{}
}

func (c *CompanyOSINT) ASNLookup(ip string) (string, error) {
url := "https://ipinfo.io/" + ip + "/json"
resp, err := http.Get(url)
if err != nil {
return "", err
}
defer resp.Body.Close()

var data map[string]interface{}
json.NewDecoder(resp.Body).Decode(&data)
if asn, ok := data["org"]; ok {
return asn.(string), nil
}
return "", nil
}

func (c *CompanyOSINT) CertTransparency(domain string) ([]string, error) {
url := "https://crt.sh/?q=" + domain + "&output=json"
resp, err := http.Get(url)
if err != nil {
return nil, err
}
defer resp.Body.Close()

var data []map[string]interface{}
json.NewDecoder(resp.Body).Decode(&data)

var certs []string
for _, item := range data {
if name, ok := item["name_value"]; ok {
certs = append(certs, name.(string))
}
}
return certs, nil
}
