module github.com/lrstanley/geoip

go 1.27.4

// TODO: https://github.com/go-chi/httprate/pull/56
replace github.com/go-chi/httprate => github.com/lrstanley/gochi-httprate v0.0.0-20260329223651-a582c3dcf8a3

require (
	github.com/go-chi/chi/v5 v5.5.0
	github.com/go-chi/httprate v0.16.1
	github.com/lrstanley/chix/v2 v2.0.0-beta.12
	github.com/lrstanley/clix/v2 v2.0.7
	github.com/lrstanley/go-bogon v1.0.1
	github.com/lrstanley/x/sync v0.0.0-20261123121849-6e1388555e13
	github.com/oschwald/maxminddb-golang v1.13.1
)

require (
	github.com/alecthomas/kong v1.17.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.16 // indirect
	github.com/go-playground/form/v4 v4.5.1 // indirect
	github.com/go-playground/locales v0.14.2 // indirect
	github.com/go-playground/universal-translator v0.18.2 // indirect
	github.com/go-playground/validator/v10 v10.30.6 // indirect
	github.com/klauspost/cpuid/v2 v2.4.2 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/lmittmann/tint v1.2.1 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/zeebo/xxh3 v1.1.0 // indirect
	golang.org/x/crypto v0.59.0 // indirect
	golang.org/x/sys v0.50.0 // indirect
	golang.org/x/text v0.44.0 // indirect
)
