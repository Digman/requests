module github.com/Digman/requests

go 1.22.0

toolchain go1.22.5

require (
	github.com/Digman/tls-client v1.0.10
	github.com/bogdanfinn/fhttp v0.5.34
	github.com/tidwall/gjson v1.17.0
)

replace github.com/Digman/tls-client => ../tls-client

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/bogdanfinn/utls v1.6.5 // indirect
	github.com/cloudflare/circl v1.6.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/quic-go/quic-go v0.49.0 // indirect
	github.com/tam7t/hpkp v0.0.0-20160821193359-2b70b4024ed5 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.31.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)
