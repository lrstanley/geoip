// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	maxminddb "github.com/oschwald/maxminddb-golang"
)

type DB struct {
	path string
}

// Note that cache may not always be filled.
type MetaCache struct {
	sync.RWMutex
	cache *maxminddb.Metadata
}

var mcache = &MetaCache{}

func (d *DB) checkForUpdates() (needsUpdate bool, err error) {
	curSeconds := time.Now().UnixNano() / int64(time.Second)
	stat, err := os.Stat(d.path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}

		return true, err
	}

	var db *maxminddb.Reader

	db, err = maxminddb.Open(d.path)
	if err != nil {
		return true, err
	}

	mcache.Lock()
	mcache.cache = &db.Metadata
	mcache.Unlock()

	if curSeconds-(stat.ModTime().UnixNano()/int64(time.Second)) < 604800 {
		return false, nil
	}

	return true, nil
}

func (d *DB) update(url string) error {
	started := time.Now()
	debug.Printf("fetching new geoip data from: %s", url)

	// Create or truncate if already exists.
	tmpfile, err := ioutil.TempFile("", "geoipdb-")
	if err != nil {
		return fmt.Errorf("unable to create temp file: %s", err)
	}
	defer func() {
		if err = tmpfile.Close(); err != nil {
			debug.Printf("error while closing %q: %s", tmpfile.Name(), err)
		}
		debug.Printf("deleting: %q", tmpfile.Name())
		if err = os.Remove(tmpfile.Name()); err != nil {
			debug.Printf("error while removing %q: %s", tmpfile.Name(), err)
		}
	}()

	debug.Printf("streaming new database to: %q", tmpfile.Name())
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}

	if _, err = io.Copy(tmpfile, gz); err != nil {
		return err
	}

	if err = gz.Close(); err != nil {
		return err
	}

	debug.Printf("successfully downloaded and decompressed new database to %q, verifying now", tmpfile.Name())
	if _, err = tmpfile.Seek(0, 0); err != nil {
		return err
	}

	db, err := maxminddb.Open(tmpfile.Name())
	if err != nil {
		return fmt.Errorf("error while attempting to verify geoip data: %s", err)
	}

	if err = db.Verify(); err != nil {
		db.Close()
		return fmt.Errorf("error while attempting to verify geoip data: %s", err)
	}

	mcache.Lock()
	mcache.cache = &db.Metadata
	mcache.Unlock()
	db.Close()

	debug.Println("verification complete, updating active database")

	file, err := os.Create(d.path)
	if err != nil {
		return err
	}

	var written int64
	written, err = io.Copy(file, tmpfile)
	if err != nil {
		return err
	}

	debug.Printf("successfully wrote %d bytes to %q (took %s)", written, file.Name(), time.Since(started))

	return nil
}

// IPSearch is the struct->tag search query to search through the Maxmind DB.
type IPSearch struct {
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
		Lat       float64 `maxminddb:"latitude"`
		Long      float64 `maxminddb:"longitude"`
		MetroCode int     `maxminddb:"metro_code"`
		TimeZone  string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`
	Subdivisions []struct {
		Code  string            `maxminddb:"iso_code"`
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
	Traits struct {
		Proxy bool `maxminddb:"is_anonymous_proxy"`
	} `maxminddb:"traits"`
}

// AddrResult contains the geolocation and host information for an IP/host.
type AddrResult struct {
	IP            net.IP   `json:"ip"`
	Summary       string   `json:"summary"`
	City          string   `json:"city"`
	Subdivision   string   `json:"subdivision"`
	Country       string   `json:"country"`
	CountryCode   string   `json:"country_abbr"`
	Continent     string   `json:"continent"`
	ContinentCode string   `json:"continent_abbr"`
	Lat           float64  `json:"latitude"`
	Long          float64  `json:"longitude"`
	Timezone      string   `json:"timezone"`
	PostalCode    string   `json:"postal_code"`
	Proxy         bool     `json:"proxy"`
	Hosts         []string `json:"hosts"`
	Error         string   `json:"error,omitempty"`
}

// addrLookup does a geoip lookup of an IP address. filters is passed into
// this function, in case there are any long running tasks which the user
// may not even want (e.g. reverse dns lookups).
func addrLookup(ctx context.Context, addr net.IP, filters []string) (*AddrResult, error) {
	var result *AddrResult
	var err error

	db, err := maxminddb.Open(flags.DBPath)
	if err != nil {
		return nil, err
	}

	var query IPSearch

	err = db.Lookup(addr, &query)
	db.Close()

	if err != nil {
		return nil, err
	}

	result = &AddrResult{
		IP:            addr,
		City:          query.City.Names["en"],
		Country:       query.Country.Names["en"],
		CountryCode:   query.Country.Code,
		Continent:     query.Continent.Names["en"],
		ContinentCode: query.Continent.Code,
		Lat:           query.Location.Lat,
		Long:          query.Location.Long,
		Timezone:      query.Location.TimeZone,
		PostalCode:    query.Postal.Code,
		Proxy:         query.Traits.Proxy,
		Hosts:         []string{},
	}

	var subdiv []string
	for i := 0; i < len(query.Subdivisions); i++ {
		subdiv = append(subdiv, query.Subdivisions[i].Names["en"])
	}
	result.Subdivision = strings.Join(subdiv, ", ")

	var summary []string
	if result.City != "" {
		summary = append(summary, result.City)
	}
	if result.Subdivision != "" {
		summary = append(summary, result.Subdivision)
	}
	if result.CountryCode != "" {
		summary = append(summary, result.CountryCode)
	} else if result.Country != "" {
		summary = append(summary, result.Country)
	} else if result.ContinentCode != "" {
		summary = append(summary, result.ContinentCode)
	} else if result.Continent != "" {
		summary = append(summary, result.Continent)
	}
	result.Summary = strings.Join(summary, ", ")

	if result.Summary == "" {
		result.Error = "no results found"
	}

	wantsHosts := len(filters) == 0
	if !wantsHosts {
		for i := 0; i < len(filters); i++ {
			if filters[i] == "hosts" {
				wantsHosts = true
				break
			}
		}
	}

	if wantsHosts {
		var names []string
		resolver := &net.Resolver{}
		dnsCtx, cancel := context.WithTimeout(ctx, flags.DNSTimeout)
		defer cancel()

		if names, err = resolver.LookupAddr(dnsCtx, addr.String()); err == nil {
			for i := 0; i < len(names); i++ {
				// These are FQDN's where absolute hosts contain a suffixed ".".
				result.Hosts = append(result.Hosts, strings.TrimSuffix(names[i], "."))
			}
		}
	}

	return result, nil
}
