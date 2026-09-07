package cred_crack

import (
"os"
"strconv"
"strings"
)

type WordlistManager struct {
BaseWordlist string
OutputPath   string
}

func NewWordlistManager(baseWordlist, outputPath string) *WordlistManager {
return &WordlistManager{
BaseWordlist: baseWordlist,
OutputPath:   outputPath,
}
}

func (w *WordlistManager) ReadWordlist() ([]string, error) {
file, err := os.Open(w.BaseWordlist)
if err != nil {
return nil, err
}
defer file.Close()
buf := make([]byte, 1024*1024)
n, err := file.Read(buf)
if err != nil {
return nil, err
}
return strings.Split(string(buf[:n]), "\n"), nil
}

func (w *WordlistManager) MutateWordlist() ([]string, error) {
words, err := w.ReadWordlist()
if err != nil {
return nil, err
}
var mutated []string
for _, word := range words {
word = strings.TrimSpace(word)
if word == "" {
continue
}
mutated = append(mutated, word)
mutated = append(mutated, strings.Title(word))
mutated = append(mutated, strings.ToUpper(word))
mutated = append(mutated, word+"1")
mutated = append(mutated, word+"123")
mutated = append(mutated, word+"!")
mutated = append(mutated, word+"2023")
mutated = append(mutated, word+"2024")
mutated = append(mutated, strings.Title(word)+"123")
mutated = append(mutated, strings.Title(word)+"2023!")
}
return mutated, nil
}

func (w *WordlistManager) WriteWordlist(words []string) bool {
file, err := os.Create(w.OutputPath)
if err != nil {
return false
}
defer file.Close()
file.WriteString(strings.Join(words, "\n"))
return true
}

func (w *WordlistManager) GenerateMutatedWordlist() bool {
mutated, err := w.MutateWordlist()
if err != nil {
return false
}
return w.WriteWordlist(mutated)
}

func (w *WordlistManager) GetWordlistSize() int {
words, _ := w.ReadWordlist()
return len(words)
}

func (w *WordlistManager) GetMutatedWordlistSize() int {
mutated, _ := w.MutateWordlist()
return len(mutated)
}

func (w *WordlistManager) CountWords() string {
words, _ := w.ReadWordlist()
return strconv.Itoa(len(words))
}
