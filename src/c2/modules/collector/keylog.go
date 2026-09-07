package collector

import (
"os"
"time"
)

type Keylogger struct {
LogFile string
}

func NewKeylogger() *Keylogger {
return &Keylogger{
LogFile: "keylog_" + time.Now().Format("20060102_150405") + ".txt",
}
}

func (k *Keylogger) Start() error {
return os.WriteFile(k.LogFile, []byte("Keylogger started\n"), 0644)
}

func (k *Keylogger) LogKey(key string) error {
f, err := os.OpenFile(k.LogFile, os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
return err
}
defer f.Close()
_, err = f.WriteString(key)
return err
}

func (k *Keylogger) Stop() error {
return os.WriteFile(k.LogFile, []byte("\nKeylogger stopped\n"), 0644)
}
