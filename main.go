// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bluele/gcache"
	gflags "github.com/jessevdk/go-flags"
)

// Should be automatically added by goreleaser.
var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

type Flags struct {
	Debug          bool          `short:"d" long:"debug" description:"enable exception display and pprof endpoints (warn: dangerous)"`
	Quiet          bool          `short:"q" long:"quiet" description:"disable verbose output"`
	DBPath         string        `long:"db" description:"path to read/store Maxmind DB" default:"geoip.db"`
	UpdateInterval time.Duration `long:"interval" description:"interval of time between database update checks" default:"2h"`
	UpdateURL      string        `long:"update-url" description:"maxmind database file download location (must be gzipped)" default:"http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"`
	Cache          struct {
		Size   int           `long:"size" description:"total number of lookups to keep in ARC cache (50% most recent, 50% most requested)" default:"500"`
		Expire time.Duration `long:"expire" description:"expiration time of cache" default:"20m"`
	} `group:"Cache Options" namespace:"cache"`
	HTTP struct {
		Bind  string   `short:"b" long:"bind" description:"address and port to bind to" default:":8080"`
		Proxy bool     `long:"proxy" description:"obey X-Forwarded-For headers (warn: dangerous, make sure to only bind to localhost)"`
		Limit int      `long:"limit" description:"number of requests/ip/hour" default:"2000"`
		CORS  []string `long:"cors" description:"cors origin domain to allow with https?:// prefix (empty => '*'; use flag multiple times)"`
		TLS   struct {
			Use  bool   `long:"use" description:"enable tls"`
			Cert string `long:"cert" description:"path to ssl certificate"`
			Key  string `long:"key" description:"path to ssl key"`
		} `group:"TLS Options" namespace:"tls"`
	} `group:"HTTP Options" namespace:"http"`
	DNS struct {
		Timeout   time.Duration `long:"timeout" description:"max allowed duration when looking up hostnames (may cause queries to be slow)" default:"2s"`
		Resolvers []string      `long:"resolver" description:"resolver (in host:port form) to use for dns lookups (doesn't work with windows and plan9) (can be used multiple times)"`
		Local     bool          `long:"uselocal" description:"adds local (system) resolvers to the list of resolvers to use"`
	} `group:"DNS Lookup Options" namespace:"dns"`
	Version bool `short:"v" long:"version" description:"print the version and compilation date"`
}

var (
	flags    Flags
	logger   = log.New(ioutil.Discard, "", log.LstdFlags|log.Lshortfile)
	db       *DB
	arc      gcache.Cache
	resolver *net.Resolver
)

func main() {
	parser := gflags.NewParser(&flags, gflags.HelpFlag)
	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if flags.Version {
		fmt.Printf("geoip version %q (compiled: %q, commit: %q)\n", version, date, commit)
		os.Exit(0)
	}

	if !flags.Quiet {
		logger.SetOutput(os.Stdout)
	}

	db = &DB{path: flags.DBPath}
	arc = gcache.New(flags.Cache.Size).ARC().Expiration(flags.Cache.Expire).Build()

	if len(flags.DNS.Resolvers) == 0 {
		resolver = net.DefaultResolver
	} else {
		resolver = &net.Resolver{PreferGo: true, Dial: customResolver}
	}

	go func() {
		var needsUpdate bool
		var err error
		for {
			logger.Println("checking for database updates")
			needsUpdate, err = db.checkForUpdates()
			if needsUpdate {
				if err != nil {
					logger.Printf("database needs update due to error (%s)", err)
				} else {
					logger.Println("database needs update")
				}

				if err = db.update(flags.UpdateURL); err != nil {
					logger.Println(err)
				}
			} else {
				logger.Println("no database updates needed")
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

func customResolver(ctx context.Context, network, address string) (net.Conn, error) {
	var index int

	if flags.DNS.Local {
		index = rand.Intn(len(flags.DNS.Resolvers) + 1)
	} else {
		// Generate a random number, which is used to select a resolver.
		// However, if the number generated is out of the bounds of the
		// amount of resolvers, use the system resolver, since they
		// requested it.
		index = rand.Intn(len(flags.DNS.Resolvers))
	}

	if index == len(flags.DNS.Resolvers) {
		return net.Dial(network, address)
	}

	addr := flags.DNS.Resolvers[index]

	if strings.Contains(addr, ":") {
		return net.Dial(network, addr)
	}
	return net.Dial(network, addr+":53")
}
