// Copyright (c) Liam Stanley <liam@liam.sh>. All rights reserved. Use of
// this source code is governed by the MIT license that can be found in
// the LICENSE file.

package cache

import (
	"time"

	gcache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lfu"
	"github.com/lrstanley/geoip/internal/metrics"
)

// New returns a new ARC cache with the given size and expiration, with a custom key
// and value type.
func New[K comparable, V any](name string, size int, expiration time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		name: name,
		gc:   gcache.New(gcache.AsLFU[K, V](lfu.WithCapacity(size))),
	}
}

type Cache[K comparable, V any] struct {
	name string
	exp  time.Duration
	gc   *gcache.Cache[K, V]
}

// Get returns the value for the given key, if it exists. Otherwise, it returns the
// default value of the value type.
func (c *Cache[K, V]) Get(key K) (val V) {
	metrics.CacheRequestCount.WithLabelValues(c.name).Inc()
	var ok bool

	val, ok = c.gc.Get(key)
	if !ok {
		return val
	}

	metrics.CacheHitCount.WithLabelValues(c.name).Inc()
	return val
}

// Set sets the value for the given key.
func (c *Cache[K, V]) Set(key K, val V) {
	c.gc.Set(key, val, gcache.WithExpiration(c.exp))
}
