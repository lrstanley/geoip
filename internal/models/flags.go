// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import "time"

type Flags struct {
	HTTP ConfigHTTP `embed:"" prefix:"http." envprefix:"HTTP_" group:"HTTP Server options"`
	DB   ConfigDB   `embed:"" prefix:"db." envprefix:"DB_" group:"DB Options"`
	DNS  ConfigDNS  `embed:"" prefix:"dns." envprefix:"DNS_" group:"DNS Lookup Options"`
}

type ConfigHTTP struct {
	BindAddr       string   `env:"BIND_ADDR" name:"bind-addr" default:":8080" required:"" help:"ip:port pair to bind to"`
	TrustedProxies []string `env:"TRUSTED_PROXIES" name:"trusted-proxies" help:"CIDR ranges that we trust the X-Forwarded-For header from (addl opts: local, *, cloudflare, and/or custom header to use)"`
	MaxConcurrent  int      `env:"MAX_CONCURRENT" name:"max-concurrent" help:"limit total max concurrent requests across all connections (0 for no limit)"`
	Limit          int      `env:"LIMIT" name:"limit" help:"number of requests/ip/hour" default:"2000"`
	HSTS           bool     `env:"HSTS" name:"hsts" help:"enable HTTP Strict Transport Security"`
	CORS           []string `env:"CORS" name:"cors" default:"*" help:"CORS allowed origins"`
}

type ConfigDB struct {
	GeoIPPath string `env:"GEOIP_PATH" name:"geoip-path" help:"path to read/store GeoIP Maxmind DB" default:"geoip.db"`
	GeoIPURL  string `env:"GEOIP_UPDATE_URL" name:"geoip-update-url" help:"GeoIP database file download location (must be gzipped)" default:"https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz"`

	ASNPath string `env:"ASN_PATH" name:"asn-path" help:"path to read/store ASN Maxmind DB" default:"asn.db"`
	ASNURL  string `env:"ASN_UPDATE_URL" name:"asn-update-url" help:"ASN database file download location (must be gzipped)" default:"https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-ASN&license_key=%s&suffix=tar.gz"`

	LicenseKey      string        `env:"LICENSE_KEY" name:"license-key" help:"maxmind license key (must register for a maxmind account)" required:""`
	UpdateInterval  time.Duration `env:"UPDATE_INTERVAL" name:"update-interval" help:"interval of time between database update checks" default:"12h"`
	CacheSize       int           `env:"CACHE_SIZE" name:"size" help:"total number of lookups to keep in ARC cache (50% most recent, 50% most requested)" default:"1000"`
	CacheExpire     time.Duration `env:"CACHE_EXPIRE" name:"expire" help:"expiration time of cache" default:"1h"`
	DefaultLanguage string        `env:"DEFAULT_LANGUAGE" name:"lang" help:"default language to use for geolocation" default:"en"`
}

type ConfigDNS struct {
	Resolvers   []string      `env:"RESOLVERS" name:"resolver" help:"resolver (in host:port form) to use for dns lookups (doesn't work with windows and plan9) (can be used multiple times)"`
	Local       bool          `env:"LOCAL" name:"uselocal" help:"adds local (system) resolvers to the list of resolvers to use"`
	CacheSize   int           `env:"CACHE_SIZE" name:"size" help:"total number of lookups to keep in ARC cache (50% most recent, 50% most requested)" default:"500"`
	CacheExpire time.Duration `env:"CACHE_EXPIRE" name:"expire" help:"expiration time of cache" default:"1h"`
	Timeout     time.Duration `env:"TIMEOUT" name:"timeout" help:"timeout for dns lookups (longer = better results but longer request duration)" default:"4s"`
}
