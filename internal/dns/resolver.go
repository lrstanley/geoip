// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package dns

import (
	"context"
	"math/rand/v2"
	"net"
	"strings"
	"time"

	"github.com/lrstanley/geoip/internal/models"
	"github.com/lrstanley/x/sync/cache"
	"github.com/lrstanley/x/sync/cache/policy/lfu"
)

// Resolver is a wrapper around a net.Resolver that caches lookups, and can be
// configured to use custom resolvers.
type Resolver struct {
	config models.ConfigDNS

	rdnsCache *cache.Cache[string, string]
	hostCache *cache.Cache[string, net.IP]

	rslv *net.Resolver
}

// NewResolver creates a new resolver, potentially using a custom resolver, if
// the config has been set. Lookups are cached.
func NewResolver(ctx context.Context, config models.ConfigDNS) *Resolver {
	r := &Resolver{
		config: config,
		rdnsCache: cache.New(
			ctx,
			cache.WithJanitorInterval[string, string](1*time.Minute),
			cache.WithLFU[string, string](lfu.WithCapacity(config.CacheSize)),
			cache.WithDefaultEntryOptions[string, string](cache.WithExpiration(config.CacheExpire)),
		),
		hostCache: cache.New(
			ctx,
			cache.WithJanitorInterval[string, net.IP](1*time.Minute),
			cache.WithLFU[string, net.IP](lfu.WithCapacity(config.CacheSize)),
			cache.WithDefaultEntryOptions[string, net.IP](cache.WithExpiration(config.CacheExpire)),
		),
	}

	if len(config.Resolvers) > 0 {
		r.rslv = &net.Resolver{PreferGo: true, Dial: newCustomResolver(config)}
	} else {
		r.rslv = net.DefaultResolver
	}

	return r
}

// GetIP does a dns lookup of a hostname, caching if successful.
func (r *Resolver) GetIP(ctx context.Context, host string) (net.IP, error) {
	ip, ok := r.hostCache.Get(host)
	if ok {
		return ip, nil
	}

	dctx, cancel := context.WithTimeout(ctx, r.config.Timeout)
	defer cancel()

	ips, err := r.rslv.LookupHost(dctx, host)
	if err != nil || len(ips) == 0 {
		return ip, err
	}

	ip = net.ParseIP(ips[0])
	r.hostCache.Set(host, ip)
	return ip, nil
}

// GetReverse does a reverse dns lookup of an IP address, caching if successful.
func (r *Resolver) GetReverse(ctx context.Context, ip net.IP) (host string, err error) {
	if v, ok := r.rdnsCache.Get(ip.String()); ok {
		return v, nil
	}

	dctx, cancel := context.WithTimeout(ctx, r.config.Timeout)
	defer cancel()

	var names []string

	if names, err = r.rslv.LookupAddr(dctx, ip.String()); err == nil && len(names) > 0 {
		host = strings.TrimSuffix(names[0], ".")
		r.rdnsCache.Set(ip.String(), host)
		return host, nil
	}

	return "", err
}

type customResolver func(ctx context.Context, network, address string) (net.Conn, error)

func newCustomResolver(config models.ConfigDNS) customResolver {
	var d net.Dialer

	return func(ctx context.Context, network, address string) (net.Conn, error) {
		var index int

		if config.Local {
			index = rand.IntN(len(config.Resolvers) + 1) //nolint:gosec // Not for security purposes.
		} else {
			// Generate a random number, which is used to select a resolver.
			// However, if the number generated is out of the bounds of the
			// amount of resolvers, use the system resolver, since they
			// requested it.
			index = rand.IntN(len(config.Resolvers)) //nolint:gosec // Not for security purposes.
		}

		if index == len(config.Resolvers) {
			return d.DialContext(ctx, network, address)
		}

		addr := config.Resolvers[index]

		if strings.Contains(addr, ":") {
			return d.DialContext(ctx, network, addr)
		}
		return d.DialContext(ctx, network, addr+":53")
	}
}
