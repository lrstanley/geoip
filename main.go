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
	Debug          bool          `env:"DEBUG" short:"d" long:"debug" description:"enable exception display and pprof endpoints (warn: dangerous)"`
	Quiet          bool          `env:"QUIET" short:"q" long:"quiet" description:"disable verbose output"`
	DBPath         string        `env:"DB_PATH" long:"db" description:"path to read/store Maxmind DB" default:"geoip.db"`
	UpdateInterval time.Duration `env:"UPDATE_INTERVAL" long:"interval" description:"interval of time between database update checks" default:"12h"`
	UpdateURL      string        `env:"MAXMIND_UPDATE_URL" long:"update-url" description:"maxmind database file download location (must be gzipped)" default:"https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz"`
	LicenseKey     string        `env:"MAXMIND_LICENSE_KEY" long:"license-key" description:"maxmind license key (must register for a maxmind account)" required:"true"`
	Cache          struct {
		Size   int           `env:"CACHE_SIZE" long:"size" description:"total number of lookups to keep in ARC cache (50% most recent, 50% most requested)" default:"500"`
		Expire time.Duration `env:"CACHE_EXPIRE" long:"expire" description:"expiration time of cache" default:"20m"`
	} `group:"Cache Options" namespace:"cache"`
	HTTP struct {
		Bind     string   `env:"HTTP_BIND" short:"b" long:"bind" description:"address and port to bind to" default:":8080"`
		Proxy    bool     `env:"HTTP_BEHIND_PROXY" long:"proxy" description:"obey X-Forwarded-For headers (warn: dangerous, make sure to only bind to localhost)"`
		Throttle int      `env:"HTTP_THROTTLE" long:"throttle" description:"limit total max concurrent requests across all connections"`
		Limit    int      `env:"HTTP_LIMIT" long:"limit" description:"number of requests/ip/hour" default:"2000"`
		CORS     []string `env:"HTTP_CORS" long:"cors" description:"cors origin domain to allow with https?:// prefix (empty => '*'; use flag multiple times)"`
		TLS      struct {
			Use  bool   `env:"TLS_USE" long:"use" description:"enable tls"`
			Cert string `env:"TLS_CERT" long:"cert" description:"path to ssl certificate"`
			Key  string `env:"TLS_KEY" long:"key" description:"path to ssl key"`
		} `group:"TLS Options" namespace:"tls"`
	} `group:"HTTP Options" namespace:"http"`
	DNS struct {
		Timeout   time.Duration `env:"DNS_TIMEOUT" long:"timeout" description:"max allowed duration when looking up hostnames (may cause queries to be slow)" default:"2s"`
		Resolvers []string      `env:"DNS_RESOLVERS" long:"resolver" description:"resolver (in host:port form) to use for dns lookups (doesn't work with windows and plan9) (can be used multiple times)"`
		Local     bool          `env:"DNS_LOCAL" long:"uselocal" description:"adds local (system) resolvers to the list of resolvers to use"`
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

				if err = db.update(flags.UpdateURL, flags.LicenseKey); err != nil {
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
