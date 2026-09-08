package cleanup

import (
"crypto/sha256"
"encoding/hex"
"os"
)

type ManifestEntry struct {
Path   string
Hash   string
}

type Manifest struct {
Entries []ManifestEntry
}

func NewManifest() *Manifest {
return &Manifest{Entries: []ManifestEntry{}}
}

func (m *Manifest) AddEntry(path string) ManifestEntry {
content, _ := os.ReadFile(path)
hash := sha256.Sum256(content)
entry := ManifestEntry{Path: path, Hash: hex.EncodeToString(hash[:])}
m.Entries = append(m.Entries, entry)
return entry
}

func (m *Manifest) ExportManifest() ManifestEntry {
return ManifestEntry{Path: "/tmp/angel_manifest", Hash: "exported"}
}
