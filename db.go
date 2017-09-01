package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	maxminddb "github.com/oschwald/maxminddb-golang"
)

type DB struct {
	path string
}

func (d *DB) checkForUpdates() (needsUpdate bool, err error) {
	curSeconds := time.Now().UnixNano() / int64(time.Second)
	stat, err := os.Stat(d.path)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}

		return true, err
	}

	_, err = maxminddb.Open(d.path)
	if err != nil {
		return true, err
	}

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
		Lat      float64 `maxminddb:"latitude"`
		Long     float64 `maxminddb:"longitude"`
		TimeZone string  `maxminddb:"time_zone"`
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

	// RegisteredCountry struct {
	// 	Code  string            `maxminddb:"iso_code"`
	// 	Names map[string]string `maxminddb:"names"`
	// } `maxminddb:"registered_country"`
	// RepresentedCountry struct {
	// 	Code  string            `maxminddb:"iso_code"`
	// 	Names map[string]string `maxminddb:"names"`
	// 	Type  string            `maxminddb:"type"`
	// } `maxminddb:"represented_country"`
}

// AddrResult contains the geolocation and host information for an IP/host.
type AddrResult struct {
	City          string
	Subdivision   string
	Country       string
	CountryCode   string
	Continent     string
	ContinentCode string
	Lat           float64
	Long          float64
	Timezone      string
	PostalCode    string
	Proxy         bool
	Hosts         []string
}

// addrLookup does a geoip lookup of an IP address
func addrLookup(path string, addr net.IP) (*AddrResult, error) {
	db, err := maxminddb.Open(path)
	if err != nil {
		return nil, err
	}

	var results IPSearch

	err = db.Lookup(addr, &results)
	db.Close()

	if err != nil {
		return nil, err
	}

	res := &AddrResult{
		City:          results.City.Names["en"],
		Country:       results.Country.Names["en"],
		CountryCode:   results.Country.Code,
		Continent:     results.Continent.Names["en"],
		ContinentCode: results.Continent.Code,
		Lat:           results.Location.Lat,
		Long:          results.Location.Long,
		Timezone:      results.Location.TimeZone,
		PostalCode:    results.Postal.Code,
		Proxy:         results.Traits.Proxy,
	}

	var subdiv []string
	for i := 0; i < len(results.Subdivisions); i++ {
		subdiv = append(subdiv, results.Subdivisions[i].Names["en"])
	}
	res.Subdivision = strings.Join(subdiv, ", ")

	if names, err := net.LookupAddr(addr.String()); err == nil {
		for i := 0; i < len(names); i++ {
			// These are FQDN's where absolute hosts contain a suffixed ".".
			res.Hosts = append(res.Hosts, strings.TrimSuffix(names[i], "."))
		}
	}

	return res, nil
}
