package cleanup

import "os"

type ArtifactResult struct {
Path   string
Status string
}

type ArtifactCleanup struct {
Artifacts []string
}

func (a ArtifactCleanup) DeleteArtifacts() ArtifactResult {
for _, path := range a.Artifacts {
os.RemoveAll(path)
}
return ArtifactResult{Path: "all_artifacts", Status: "success"}
}

func (a ArtifactCleanup) DeleteLogs() ArtifactResult {
os.RemoveAll("/var/log/angel")
return ArtifactResult{Path: "/var/log/angel", Status: "success"}
}

func (a ArtifactCleanup) DeleteConfigs() ArtifactResult {
os.RemoveAll("/etc/angel")
return ArtifactResult{Path: "/etc/angel", Status: "success"}
}
