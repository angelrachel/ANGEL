package ai_orchestrator

import (
	"time"
)

type AttackChain struct {
	Stage    string
	Success  bool
	Started  time.Time
	Finished time.Time
	Router   *AIRouter
}

func NewAttackChain() *AttackChain {
	return &AttackChain{
		Router:  NewAIRouter(),
		Started: time.Now(),
	}
}

func (a *AttackChain) StartRecon(target string) {
	a.Router.SetTarget(target)
	a.Stage = "Recon"
}

func (a *AttackChain) AnalyzeRecon(openPorts []int) {
	bestVector := ""
	priority := 0

	for _, port := range openPorts {
		switch port {
		case 445:
			bestVector = "SMB Attack (Pass-the-Hash)"
			priority = 10
		case 1433:
			bestVector = "MSSQL Attack (SQL Injection)"
			priority = 9
		case 3306:
			bestVector = "MySQL Attack (SQL Injection)"
			priority = 8
		case 8080:
			bestVector = "Web Application Attack"
			priority = 7
		default:
			bestVector = "Generic Service Attack"
			priority = 1
		}
		a.Router.AssessRoute(bestVector, "Found open port "+itoa(port), priority)
	}
}

func (a *AttackChain) ExecuteBestRoute() string {
	route := a.Router.SelectBestRoute()
	return a.Router.ExecuteRoute(route)
}

func (a *AttackChain) MarkSuccess() {
	a.Success = true
	a.Finished = time.Now()
}

func (a *AttackChain) GetStage() string {
	return a.Stage
}

func (a *AttackChain) IsSuccess() bool {
	return a.Success
}

func (a *AttackChain) GetDuration() time.Duration {
	return a.Finished.Sub(a.Started)
}

func itoa(num int) string {
	if num == 0 {
		return "0"
	}
	digits := []byte{}
	for num > 0 {
		digits = append([]byte{byte('0' + num%10)}, digits...)
		num /= 10
	}
	return string(digits)
}
