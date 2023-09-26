module github.com/Digman/requests

go 1.19

require (
	github.com/Digman/tls-client v1.0.8
	github.com/bogdanfinn/fhttp v0.5.24
	github.com/tidwall/gjson v1.17.0
)

replace github.com/Digman/tls-client => ../tls-client

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/bogdanfinn/utls v1.5.16 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	golang.org/x/crypto v0.13.0 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)
