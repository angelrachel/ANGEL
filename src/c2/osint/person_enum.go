package osint

import (
"io"
"net/http"
"net/url"
"strings"
"time"
)

type PersonEnum struct {
Client *http.Client
}

func NewPersonEnum() *PersonEnum {
return &PersonEnum{
Client: &http.Client{Timeout: 10 * time.Second},
}
}

func (p *PersonEnum) SearchPerson(name string) (string, error) {
searchURL := "https://www.google.com/search?q=" + url.QueryEscape(name)
resp, err := p.Client.Get(searchURL)
if err != nil {
return "", err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil {
return "", err
}
return string(body), nil
}

func (p *PersonEnum) SearchSocialMedia(name string) (string, error) {
searchURL := "https://www.google.com/search?q=site:linkedin.com+" + url.QueryEscape(name)
resp, err := p.Client.Get(searchURL)
if err != nil {
return "", err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil {
return "", err
}
return string(body), nil
}

func (p *PersonEnum) SearchEmail(email string) (string, error) {
searchURL := "https://www.google.com/search?q=" + url.QueryEscape(email)
resp, err := p.Client.Get(searchURL)
if err != nil {
return "", err
}
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)
if err != nil {
return "", err
}
return string(body), nil
}

func (p *PersonEnum) DumpAllData(person string) (string, error) {
var result string
result += "Searching person: " + person + "\n"
google, _ := p.SearchPerson(person)
result += strings.TrimSpace(google) + "\n"
linkedin, _ := p.SearchSocialMedia(person)
result += strings.TrimSpace(linkedin) + "\n"
return result, nil
}
