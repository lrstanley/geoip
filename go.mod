module github.com/lrstanley/geoip

go 1.26.0

// TODO: https://github.com/go-chi/httprate/pull/56
replace github.com/go-chi/httprate => github.com/lrstanley/gochi-httprate v0.0.0-20260329223651-a582c3dcf8a3

require (
	github.com/go-chi/chi/v5 v5.2.3
	github.com/go-chi/httprate v0.15.0
	github.com/lrstanley/chix/v2 v2.0.0-beta.0
	github.com/lrstanley/clix/v2 v2.0.0-beta.1
	github.com/lrstanley/go-bogon v1.0.0
	github.com/lrstanley/x/scheduler v0.0.0-20260329042521-db1db209ec02
	github.com/lrstanley/x/sync v0.0.0-20260329042521-db1db209ec02
	github.com/oschwald/maxminddb-golang v1.13.1
)

require (
	github.com/alecthomas/kong v1.12.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.11 // indirect
	github.com/go-playground/form/v4 v4.3.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.28.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lmittmann/tint v1.1.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
)
