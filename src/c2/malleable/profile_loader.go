package malleable

import (
"fmt"
"os"

"gopkg.in/yaml.v3"
)

type C2Profile struct {
UserAgent  string   `yaml:"user_agent"`
URI        []string `yaml:"uri"`
Headers    []string `yaml:"headers"`
Server     string   `yaml:"server"`
StatusCode int      `yaml:"status_code"`
}

func LoadProfile(path string) (*C2Profile, error) {
data, err := os.ReadFile(path)
if err != nil {
return nil, fmt.Errorf("failed to read profile: %w", err)
}

var profile C2Profile
if err := yaml.Unmarshal(data, &profile); err != nil {
return nil, fmt.Errorf("failed to parse profile: %w", err)
}

return &profile, nil
}

func (p *C2Profile) ValidateUserAgent(ua string) bool {
return ua == p.UserAgent
}

func (p *C2Profile) ValidateURI(uri string) bool {
for _, u := range p.URI {
if u == uri {
return true
}
}
return false
}
