module github.com/lrstanley/geoip

go 1.26.0

// TODO: https://github.com/go-chi/httprate/pull/56
replace github.com/go-chi/httprate => github.com/lrstanley/gochi-httprate v0.0.0-20260329223651-a582c3dcf8a3

require (
	github.com/go-chi/chi/v5 v5.2.5
	github.com/go-chi/httprate v0.15.0
	github.com/lrstanley/chix/v2 v2.0.0-beta.1
	github.com/lrstanley/clix/v2 v2.0.0
	github.com/lrstanley/go-bogon v1.0.0
	github.com/lrstanley/x/sync v0.0.0-20260331013828-98de5249208d
	github.com/oschwald/maxminddb-golang v1.13.1
)

require (
	github.com/alecthomas/kong v1.14.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/go-playground/form/v4 v4.3.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.2 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lmittmann/tint v1.1.3 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
)
