package cleanup

import (
"crypto/sha256"
"encoding/hex"
"encoding/json"
"os"
"time"
)

type DeletionManifest struct {
Timestamp   string   `json:"timestamp"`
FilesRemoved []string `json:"files_removed"`
Checksum    string   `json:"checksum"`
Verified    bool     `json:"verified"`
}

type Manifest struct {
ManifestPath string
}

func NewManifest(path string) *Manifest {
return &Manifest{
ManifestPath: path,
}
}

func (m *Manifest) Generate(files []string) (*DeletionManifest, error) {
manifest := &DeletionManifest{
Timestamp:    time.Now().UTC().Format(time.RFC3339),
FilesRemoved: files,
Verified:     false,
}

// Calculate checksum
data, _ := json.Marshal(manifest)
hash := sha256.Sum256(data)
manifest.Checksum = hex.EncodeToString(hash[:])

// Save manifest
jsonData, err := json.MarshalIndent(manifest, "", "  ")
if err != nil {
return nil, err
}
err = os.WriteFile(m.ManifestPath, jsonData, 0644)
if err != nil {
return nil, err
}
return manifest, nil
}

func (m *Manifest) Verify() (bool, error) {
data, err := os.ReadFile(m.ManifestPath)
if err != nil {
return false, err
}

var manifest DeletionManifest
err = json.Unmarshal(data, &manifest)
if err != nil {
return false, err
}

// Verify checksum
dataWithoutChecksum, _ := json.Marshal(manifest)
hash := sha256.Sum256(dataWithoutChecksum)
if hex.EncodeToString(hash[:]) != manifest.Checksum {
return false, nil
}
return true, nil
}

func (m *Manifest) Export() (string, error) {
data, err := os.ReadFile(m.ManifestPath)
if err != nil {
return "", err
}
return string(data), nil
}
