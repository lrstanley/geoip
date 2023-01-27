// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import "time"

type Flags struct {
	HTTP ConfigHTTP `group:"HTTP Server options" namespace:"http" env-namespace:"HTTP"`
	DB   ConfigDB   `group:"DB Options" namespace:"db" env-namespace:"DB"`
	DNS  ConfigDNS  `group:"DNS Lookup Options" namespace:"dns" env-namespace:"DNS"`
}

type ConfigHTTP struct {
	BindAddr       string   `env:"BIND_ADDR"       long:"bind-addr"       default:":8080" required:"true" description:"ip:port pair to bind to"`
	TrustedProxies []string `env:"TRUSTED_PROXIES" long:"trusted-proxies" env-delim:"," description:"CIDR ranges that we trust the X-Forwarded-For header from (addl opts: local, *, cloudflare, and/or custom header to use)"`
	MaxConcurrent  int      `env:"MAX_CONCURRENT"  long:"max-concurrent"  description:"limit total max concurrent requests across all connections (0 for no limit)"`
	Limit          int      `env:"LIMIT"           long:"limit"           description:"number of requests/ip/hour" default:"2000"`
	HSTS           bool     `env:"HSTS"            long:"hsts"            description:"enable HTTP Strict Transport Security"`
	CORS           []string `env:"CORS"            long:"cors"            env-delim:"," default:"*" description:"CORS allowed origins"`
	Metrics        bool     `env:"METRICS"         long:"metrics"         description:"enable prometheus metrics on /metrics to internal IPs"`
}

type ConfigDB struct {
	GeoIPPath string `env:"GEOIP_PATH"            long:"geoip-path" description:"path to read/store GeoIP Maxmind DB" default:"geoip.db"`
	GeoIPURL  string `env:"GEOIP_UPDATE_URL"      long:"geoip-update-url" description:"GeoIP database file download location (must be gzipped)" default:"https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz"`

	ASNPath string `env:"ASN_PATH"            long:"asn-path" description:"path to read/store ASN Maxmind DB" default:"asn.db"`
	ASNURL  string `env:"ASN_UPDATE_URL"      long:"asn-update-url" description:"ASN database file download location (must be gzipped)" default:"https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-ASN&license_key=%s&suffix=tar.gz"`

	LicenseKey      string        `env:"LICENSE_KEY"     long:"license-key" description:"maxmind license key (must register for a maxmind account)" required:"true"`
	UpdateInterval  time.Duration `env:"UPDATE_INTERVAL" long:"update-interval" description:"interval of time between database update checks" default:"12h"`
	CacheSize       int           `env:"CACHE_SIZE"   long:"size" description:"total number of lookups to keep in ARC cache (50% most recent, 50% most requested)" default:"1000"`
	CacheExpire     time.Duration `env:"CACHE_EXPIRE" long:"expire" description:"expiration time of cache" default:"1h"`
	DefaultLanguage string        `env:"DEFAULT_LANGUAGE" long:"lang" description:"default language to use for geolocation" default:"en"`
}

type ConfigDNS struct {
	Resolvers   []string      `env:"RESOLVERS"    long:"resolver" env-delim:"," description:"resolver (in host:port form) to use for dns lookups (doesn't work with windows and plan9) (can be used multiple times)"`
	Local       bool          `env:"LOCAL"        long:"uselocal" description:"adds local (system) resolvers to the list of resolvers to use"`
	CacheSize   int           `env:"CACHE_SIZE"   long:"size"     description:"total number of lookups to keep in ARC cache (50% most recent, 50% most requested)" default:"500"`
	CacheExpire time.Duration `env:"CACHE_EXPIRE" long:"expire"   description:"expiration time of cache" default:"1h"`
	Timeout     time.Duration `env:"TIMEOUT" long:"timeout" description:"timeout for dns lookups (longer = better results but longer request duration)" default:"4s"`
}
