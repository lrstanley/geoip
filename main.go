// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	gflags "github.com/jessevdk/go-flags"
)

// TODO:
//  - https://github.com/yl2chen/cidranger (or lrstanley/go-bogon)
//  - https://github.com/bluele/gcache
//  - Provide google map direct link in full API?

var version, commit, date = "unknown", "unknown", "unknown"

type Flags struct {
	Debug          bool          `short:"d" long:"debug" description:"enable debug output"`
	DBPath         string        `long:"db" description:"path to read/store Maxmind DB" default:"geoip.db"`
	UpdateInterval time.Duration `long:"interval" description:"interval of time between database update checks" default:"2h"`
	UpdateURL      string        `long:"update-url" description:"maxmind database file download location (must be gzipped)" default:"http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"`
	NoBogon        bool          `long:"no-bogon" description:"disallow bogon addresses in lookups"`
	HTTP           struct {
		Bind  string   `short:"b" long:"bind" description:"address and port to bind to" default:":8080"`
		Proxy bool     `long:"proxy" description:"obey X-Forwarded-For headers (dangerous!)"`
		Limit int      `long:"limit" description:"number of requests/ip/hour" default:"2000"`
		CORS  []string `long:"cors" description:"cors origin domain to allow (empty => '*'; use flag multiple times)"`
		TLS   struct {
			Use  bool   `long:"use" description:"enable tls"`
			Cert string `long:"cert" description:"path to ssl certificate"`
			Key  string `long:"key" description:"path to ssl key"`
		} `group:"TLS Options" namespace:"tls"`
	} `group:"HTTP Options" namespace:"http"`
}

var flags Flags
var debug = log.New(ioutil.Discard, "", log.LstdFlags|log.Lshortfile)
var db *DB

func main() {
	parser := gflags.NewParser(&flags, gflags.HelpFlag)
	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if flags.Debug {
		debug.SetOutput(os.Stdout)
	}

	db = &DB{path: flags.DBPath}

	go func() {
		var needsUpdate bool
		var err error
		for {
			debug.Println("checking for database updates")
			needsUpdate, err = db.checkForUpdates()
			if needsUpdate {
				if err != nil {
					debug.Printf("database needs update due to error (%s)", err)
				} else {
					debug.Println("database needs update")
				}

				if err = db.update(flags.UpdateURL); err != nil {
					debug.Println(err)
				}
			} else {
				debug.Println("no database updates needed")
			}

			time.Sleep(flags.UpdateInterval)
		}
	}()

	httpCloser := make(chan struct{})
	go initHTTP(httpCloser)

	catch()
	close(httpCloser)
	fmt.Println("exiting")
}

func catch() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	fmt.Println("listening for signal. CTRL+C to quit.")
	<-signals
	fmt.Println("\ninvoked termination, cleaning up")
}
