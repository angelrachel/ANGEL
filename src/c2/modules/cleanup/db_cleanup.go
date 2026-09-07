package cleanup

import (
"database/sql"
)

type DBCleanup struct {
DB *sql.DB
}

func NewDBCleanup(db *sql.DB) *DBCleanup {
return &DBCleanup{
DB: db,
}
}

func (d *DBCleanup) DeleteJavaObjects() error {
if d.DB == nil {
return nil
}
_, err := d.DB.Exec("DROP JAVA SOURCE khunt")
return err
}

func (d *DBCleanup) DeleteStoredProcedures() error {
if d.DB == nil {
return nil
}
procedures := []string{"khunt_cmd", "khunt_hash", "khunt_fs"}
for _, proc := range procedures {
d.DB.Exec("DROP PROCEDURE " + proc)
}
return nil
}

func (d *DBCleanup) DeleteAdminAccounts() error {
if d.DB == nil {
return nil
}
_, err := d.DB.Exec("DELETE FROM users WHERE username = 'angel_backup'")
return err
}

func (d *DBCleanup) RevertChanges() error {
if d.DB == nil {
return nil
}
// Revert any schema changes
return nil
}
