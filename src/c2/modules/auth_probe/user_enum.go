package auth_probe

import (
"bytes"
"encoding/json"
"io"
"net/http"
"net/url"
"strings"
"time"
)

type UserEnumResult struct {
Username string
Exists   bool
Response string
}

type UserEnumerator struct {
URL         string
Method      string
Client      *http.Client
SuccessCode int
SuccessWord string
}

func NewUserEnumerator(url, method string, successCode int, successWord string) *UserEnumerator {
return &UserEnumerator{
URL:         url,
Method:      method,
Client:      &http.Client{Timeout: 10 * time.Second},
SuccessCode: successCode,
SuccessWord: successWord,
}
}

func (u *UserEnumerator) Enumerate(username string) (bool, error) {
var req *http.Request
var err error

if u.Method == "POST" {
form := url.Values{}
form.Set("username", username)
form.Set("password", "invalidpassword123")
req, err = http.NewRequest("POST", u.URL, bytes.NewBufferString(form.Encode()))
req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
} else {
req, err = http.NewRequest("GET", u.URL+"?username="+url.QueryEscape(username)+"&password=invalidpassword123", nil)
}

if err != nil {
return false, err
}

resp, err := u.Client.Do(req)
if err != nil {
return false, err
}
defer resp.Body.Close()

body, _ := io.ReadAll(resp.Body)
content := string(body)

if u.SuccessWord != "" {
return strings.Contains(content, u.SuccessWord), nil
}
return resp.StatusCode == u.SuccessCode, nil
}

func (u *UserEnumerator) EnumerateList(usernames []string) []UserEnumResult {
var results []UserEnumResult
for _, username := range usernames {
exists, _ := u.Enumerate(username)
results = append(results, UserEnumResult{Username: username, Exists: exists})
}
return results
}

func (u *UserEnumerator) GetValidUsers(usernames []string) []string {
var validUsers []string
for _, result := range u.EnumerateList(usernames) {
if result.Exists {
validUsers = append(validUsers, result.Username)
}
}
return validUsers
}

func (u *UserEnumerator) DumpAll(usernames []string) (string, error) {
results := u.EnumerateList(usernames)
data, _ := json.Marshal(results)
return string(data), nil
}
