package cleanup

import "os/exec"

type DBCleanupResult struct {
	Database string
	Status   string
}

type DBCleanup struct {
	Host     string
	User     string
	Password string
}

func (d DBCleanup) DropSchema() DBCleanupResult {
	cmd := exec.Command("mysql", "-h", d.Host, "-u", d.User, "-p"+d.Password, "-e", "DROP DATABASE IF EXISTS angel;")
	cmd.Run()
	return DBCleanupResult{Database: "angel", Status: "success"}
}

func (d DBCleanup) DropStoredProc() DBCleanupResult {
	cmd := exec.Command("mysql", "-h", d.Host, "-u", d.User, "-p"+d.Password, "-e", "DROP PROCEDURE IF EXISTS angel_sp;")
	cmd.Run()
	return DBCleanupResult{Database: "angel", Status: "success"}
}

func (d DBCleanup) DeleteAdminAccounts() DBCleanupResult {
	cmd := exec.Command("mysql", "-h", d.Host, "-u", d.User, "-p"+d.Password, "-e", "DELETE FROM mysql.user WHERE user='admin';")
	cmd.Run()
	return DBCleanupResult{Database: "angel", Status: "success"}
}
