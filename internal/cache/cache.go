// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package cache

import (
	"time"

	"github.com/bluele/gcache"
)

// New returns a new ARC cache with the given size and expiration, with a custom key
// and value type.
func New[K comparable, V any](size int, expiration time.Duration) *Cache[K, V] {
	c := &Cache[K, V]{
		gc: gcache.New(size).ARC().Expiration(expiration).Build(),
	}

	return c
}

type Cache[K comparable, V any] struct {
	gc gcache.Cache
}

// Get returns the value for the given key, if it exists. Otherwise, it returns the
// default value of the value type.
func (c *Cache[K, V]) Get(key K) (val V) {
	tmp, err := c.gc.GetIFPresent(key)
	if err != nil {
		return
	}

	return tmp.(V)
}

// Set sets the value for the given key.
func (c *Cache[K, V]) Set(key K, val V) {
	_ = c.gc.Set(key, val)
}
