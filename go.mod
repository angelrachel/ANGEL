module ANGEL

go 1.26.1

replace ANGEL/src/c2/server/crypto => ./src/c2/server/crypto

require (
	github.com/hirochachacha/go-smb2 v1.1.0
	golang.org/x/sys v0.47.0
)

require (
	github.com/geoffgarside/ber v1.1.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
)
