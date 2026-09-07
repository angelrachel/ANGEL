module ANGEL

go 1.26.1

replace ANGEL/src/c2/server/crypto => ./src/c2/server/crypto

require (
	github.com/denisenkom/go-mssqldb v0.12.3
	github.com/go-sql-driver/mysql v1.10.1
	github.com/godror/godror v0.51.4
	github.com/lib/pq v1.12.3
	github.com/mattn/go-sqlite3 v1.14.52
	gopkg.in/yaml.v3 v3.0.1
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/VictoriaMetrics/easyproto v0.1.4 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/godror/knownpb v0.3.0 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6 // indirect
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
