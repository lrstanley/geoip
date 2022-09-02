// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package models

import "net"

// LookupOptions are the options for the lookup service.
type LookupOptions struct {
	Languages         []string `form:"languages,omitempty" json:"languages" validate:"omitempty,dive,bcp47_language_tag|alpha,min=2,max=5"`
	DisableHostLookup bool     `form:"disable_host_lookup,omitempty" json:"disable_host_lookup"`
}

// GeoQuery is the struct->tag search query to search through the GeoIP Maxmind DB.
type GeoQuery struct {
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Country struct {
		Code  string            `maxminddb:"iso_code"`
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	Continent struct {
		Code  string            `maxminddb:"code"`
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"continent"`
	Location struct {
		Lat              float64 `maxminddb:"latitude"`
		Long             float64 `maxminddb:"longitude"`
		MetroCode        int     `maxminddb:"metro_code"`
		TimeZone         string  `maxminddb:"time_zone"`
		AccuracyRadiusKM int     `maxminddb:"accuracy_radius"`
	} `maxminddb:"location"`
	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`
	Subdivisions []struct {
		Code  string            `maxminddb:"iso_code"`
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
}

// ASNQuery is the ASN search query to search through the ASN Maxmind DB.
type ASNQuery struct {
	AutonomousSystemNumber int    `maxminddb:"autonomous_system_number"`
	AutonomousSystemOrg    string `maxminddb:"autonomous_system_organization"`

	Network *net.IPNet
}

// Response contains the geolocation and host information for an IP/host.
type Response struct {
	Query string `json:"query"`
	Error string `json:"error,omitempty"`

	Host string `json:"host,omitempty"`

	// GeoIP information.
	IP               string  `json:"ip"`
	IPType           int     `json:"ip_type"`
	Summary          string  `json:"summary"`
	City             string  `json:"city"`
	Subdivision      string  `json:"subdivision"`
	Country          string  `json:"country"`
	CountryCode      string  `json:"country_abbr"`
	Continent        string  `json:"continent"`
	ContinentCode    string  `json:"continent_abbr"`
	Lat              float64 `json:"latitude"`
	Long             float64 `json:"longitude"`
	Timezone         string  `json:"timezone"`
	PostalCode       string  `json:"postal_code"`
	AccuracyRadiusKM int     `json:"accuracy_radius_km"`

	// ASN information.
	Network string `json:"network"`
	ASN     string `json:"asn"`
	ASNOrg  string `json:"asn_org"`
}
