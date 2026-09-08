package credential

import (
"bufio"
"os"
"strings"
)

type HistoryResult struct {
Path    string
History []string
}

func ReadBashHistory(historyPath string) HistoryResult {
result := HistoryResult{Path: historyPath, History: []string{}}
file, err := os.Open(historyPath)
if err != nil {
return result
}
defer file.Close()
scanner := bufio.NewScanner(file)
for scanner.Scan() {
line := strings.TrimSpace(scanner.Text())
if line != "" {
result.History = append(result.History, line)
}
}
return result
}

func ReadZshHistory(historyPath string) HistoryResult {
result := HistoryResult{Path: historyPath, History: []string{}}
file, err := os.Open(historyPath)
if err != nil {
return result
}
defer file.Close()
scanner := bufio.NewScanner(file)
for scanner.Scan() {
line := strings.TrimSpace(scanner.Text())
if line != "" {
result.History = append(result.History, line)
}
}
return result
}
