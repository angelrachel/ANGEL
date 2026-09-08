package persistence

import "os/exec"

type ServiceResult struct {
Name   string
Status string
}

func CreateSystemdService(unitName, command string) ServiceResult {
cmd := exec.Command("bash", "-c", "echo '[Unit]' > /etc/systemd/system/"+unitName+".service && echo '[Service]' >> /etc/systemd/system/"+unitName+".service && echo 'ExecStart="+command+"' >> /etc/systemd/system/"+unitName+".service && echo '[Install]' >> /etc/systemd/system/"+unitName+".service && echo 'WantedBy=multi-user.target' >> /etc/systemd/system/"+unitName+".service")
cmd.Run()
return ServiceResult{Name: unitName, Status: "success"}
}

func StartSystemdService(unitName string) ServiceResult {
cmd := exec.Command("systemctl", "start", unitName)
cmd.Run()
return ServiceResult{Name: unitName, Status: "success"}
}

func EnableSystemdService(unitName string) ServiceResult {
cmd := exec.Command("systemctl", "enable", unitName)
cmd.Run()
return ServiceResult{Name: unitName, Status: "success"}
}
