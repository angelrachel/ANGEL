package cleanup

import (
"os"
"path/filepath"
)

type ArtifactCleanup struct{}

func NewArtifactCleanup() *ArtifactCleanup {
return &ArtifactCleanup{}
}

func (a *ArtifactCleanup) DeleteTools() error {
toolPaths := []string{
"/tmp/angel_tools",
"/opt/angel",
"/usr/local/bin/angel",
}
for _, path := range toolPaths {
os.RemoveAll(path)
}
return nil
}

func (a *ArtifactCleanup) DeleteLogs() error {
logPaths := []string{
"/var/log/angel/",
"/tmp/angel_logs",
}
for _, path := range logPaths {
os.RemoveAll(path)
}
return nil
}

func (a *ArtifactCleanup) DeleteConfigs() error {
configPaths := []string{
"/etc/angel/",
"/tmp/angel_configs",
}
for _, path := range configPaths {
os.RemoveAll(path)
}
return nil
}

func (a *ArtifactCleanup) DeleteBackups() error {
backupPaths := []string{
"/tmp/angel_backups",
"/var/backups/angel/",
}
for _, path := range backupPaths {
os.RemoveAll(path)
}
return nil
}

func (a *ArtifactCleanup) CleanupAll() error {
a.DeleteTools()
a.DeleteLogs()
a.DeleteConfigs()
a.DeleteBackups()
return nil
}
