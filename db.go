// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
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
type metaCache struct {
	sync.RWMutex
	cache *maxminddb.Metadata
}

var mcache = &metaCache{}

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

func (d *DB) update(url, licenseKey string) error {
	started := time.Now()
	url = fmt.Sprintf(url, licenseKey)

	logger.Printf("fetching new geoip data from: %s", url)

	archiveTempFile, err := ioutil.TempFile("", "geoip-archive-")
	if err != nil {
		return fmt.Errorf("unable to create temp file: %s", err)
	}
	defer func() {
		if err = archiveTempFile.Close(); err != nil {
			logger.Printf("error while closing %q: %s", archiveTempFile.Name(), err)
		}
		logger.Printf("deleting: %q", archiveTempFile.Name())
		if err = os.Remove(archiveTempFile.Name()); err != nil {
			logger.Printf("error while removing %q: %s", archiveTempFile.Name(), err)
		}
	}()

	dbTempFile, err := ioutil.TempFile("", "geoip-db-")
	if err != nil {
		return fmt.Errorf("unable to create temp file: %s", err)
	}
	defer func() {
		if err = dbTempFile.Close(); err != nil {
			logger.Printf("error while closing %q: %s", dbTempFile.Name(), err)
		}
		logger.Printf("deleting: %q", dbTempFile.Name())
		if err = os.Remove(dbTempFile.Name()); err != nil {
			logger.Printf("error while removing %q: %s", dbTempFile.Name(), err)
		}
	}()

	logger.Printf("streaming new database archive to: %q", dbTempFile.Name())
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

	tarReader := tar.NewReader(gz)
	dbFound := false

	// loop through the tar file and extract first .mmdb file we find in the
	// archive.
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("unable to read tar archive: %w", err)
		}

		if header.Typeflag == tar.TypeReg && strings.HasSuffix(header.Name, ".mmdb") {
			logger.Printf("found database in tar archive, extracting and writing to file: %v", dbTempFile.Name())
			if _, err := io.Copy(dbTempFile, tarReader); err != nil {
				return fmt.Errorf("error extracting database from tar archive: %w", err)
			}
			dbFound = true
		}
	}

	if !dbFound {
		return errors.New("no database file found in tar archive")
	}

	if err := gz.Close(); err != nil {
		return err
	}

	logger.Printf("successfully downloaded and decompressed new database to %q, verifying now", dbTempFile.Name())
	if _, err = dbTempFile.Seek(0, 0); err != nil {
		return err
	}

	db, err := maxminddb.Open(dbTempFile.Name())
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

	logger.Println("verification complete, updating active database")

	file, err := os.Create(d.path)
	if err != nil {
		return err
	}

	var written int64
	written, err = io.Copy(file, dbTempFile)
	if err != nil {
		return err
	}

	logger.Printf("successfully wrote %d bytes to %q (took %s)", written, file.Name(), time.Since(started))

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
	IP            net.IP  `json:"ip"`
	Summary       string  `json:"summary"`
	City          string  `json:"city"`
	Subdivision   string  `json:"subdivision"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_abbr"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continent_abbr"`
	Lat           float64 `json:"latitude"`
	Long          float64 `json:"longitude"`
	Timezone      string  `json:"timezone"`
	PostalCode    string  `json:"postal_code"`
	Proxy         bool    `json:"proxy"`
	Host          string  `json:"host"`
	Error         string  `json:"error,omitempty"`
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

	if result.Subdivision != "" && result.City != result.Subdivision {
		summary = append(summary, result.Subdivision)
	}

	if result.Country != "" && len(summary) == 0 {
		summary = append(summary, result.Country)
	} else if result.CountryCode != "" {
		summary = append(summary, result.CountryCode)
	}

	if result.Continent != "" && len(summary) == 0 {
		summary = append(summary, result.Continent)
	} else if result.ContinentCode != "" && result.Subdivision == "" && result.City == "" {
		summary = append(summary, result.ContinentCode)
	}

	result.Summary = strings.Join(summary, ", ")

	if result.Summary == "" {
		result.Error = "no results found"
	}

	wantsHosts := len(filters) == 0
	if !wantsHosts {
		for i := 0; i < len(filters); i++ {
			if filters[i] == "host" {
				wantsHosts = true
				break
			}
		}
	}

	if wantsHosts {
		result.Host, _ = lookupHost(ctx, addr)
	}

	return result, nil
}

func lookupHost(ctx context.Context, addr net.IP) (string, error) {
	dnsCtx, cancel := context.WithTimeout(ctx, flags.DNS.Timeout)
	defer cancel()

	var names []string
	var err error

	if names, err = resolver.LookupAddr(dnsCtx, addr.String()); err == nil && len(names) > 0 {
		return strings.TrimSuffix(names[0], "."), nil
	}

	return "", err
}
