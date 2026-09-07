package destruction

import (
"database/sql"
"fmt"
)

type DatabaseDestruction struct {
DB *sql.DB
}

func NewDatabaseDestruction(db *sql.DB) *DatabaseDestruction {
return &DatabaseDestruction{DB: db}
}

func (d *DatabaseDestruction) DropSchemas() error {
rows, err := d.DB.Query("SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('information_schema', 'pg_catalog')")
if err != nil {
return err
}
defer rows.Close()

for rows.Next() {
var schema string
if err := rows.Scan(&schema); err != nil {
continue
}
_, err := d.DB.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema))
if err != nil {
return err
}
}
return nil
}

func (d *DatabaseDestruction) DropTables() error {
rows, err := d.DB.Query("SELECT table_name, table_schema FROM information_schema.tables WHERE table_schema NOT IN ('information_schema', 'pg_catalog')")
if err != nil {
return err
}
defer rows.Close()

for rows.Next() {
var table, schema string
if err := rows.Scan(&table, &schema); err != nil {
continue
}
_, err := d.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s.%s CASCADE", schema, table))
if err != nil {
return err
}
}
return nil
}

func (d *DatabaseDestruction) DisableForeignKeys() error {
_, err := d.DB.Exec("SET session_replication_role = 'replica'")
return err
}

func (d *DatabaseDestruction) DeleteBackups() error {
_, err := d.DB.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = 'template1'")
return err
}
