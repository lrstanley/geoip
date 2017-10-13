// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.
//
// Parts of this file are copyright the github.com/go-web/httprl authors.
// see the following link for details:
//  - https://github.com/go-web/httprl/blob/master/LICENSE

package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-web/httprl"
)

// httprl's interface{} implementation currently has no way of obtaining the
// current rate limit without having the check itself count against the
// connections total limit. As such, this will have to be done manually.
func rateHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if flags.HTTP.Limit <= 0 {
			next.ServeHTTP(w, r)
			return
		}

		rate, remttl := mapLimiter.Get(httprl.DefaultKeyMaker(r), 60*60)
		remaining := uint64(flags.HTTP.Limit) - rate
		if remaining < 0 {
			remaining = 0
		}

		w.Header().Set("X-Ratelimit-Limit", fmt.Sprintf("%d", flags.HTTP.Limit))
		w.Header().Set("X-Ratelimit-Remaining", fmt.Sprintf("%d", remaining))
		w.Header().Set("X-Ratelimit-Reset", fmt.Sprintf("%d", remttl))

		next.ServeHTTP(w, r)
	})
}

// MapLimiter is a rate limiter implementation for github.com/go-web/httprl
// which is like the builtin Map limiter, but allows querying the current
// limit and expiration time.
type MapLimiter struct {
	m    sync.Mutex
	s    map[string]*rldata
	p    time.Duration
	stop chan struct{}
}

type rldata struct {
	Count  uint64
	Expire time.Time
}

// NewMapLimiter creates and initializes a new MapLimiter. The precision
// determines how often the map is scanned for expired keys, in seconds.
func NewMapLimiter(precision int32) *MapLimiter {
	return &MapLimiter{
		s: make(map[string]*rldata),
		p: time.Duration(precision) * time.Second,
	}
}

// Get returns the current amount of tracked hits for the given key.
func (m *MapLimiter) Get(key string, ttlsec int32) (count uint64, remttl int32) {
	m.m.Lock()
	defer m.m.Unlock()
	v, ok := m.s[key]
	if !ok {
		return 0, ttlsec
	}

	rttl := v.Expire.Sub(time.Now()).Seconds()
	if rttl < 1 {
		return v.Count, 0
	}

	return v.Count, int32(rttl)
}

// Hit implements the httprl.Backend interface.
func (m *MapLimiter) Hit(key string, ttlsec int32) (count uint64, remttl int32, err error) {
	m.m.Lock()
	defer m.m.Unlock()
	v, ok := m.s[key]
	if !ok {
		m.s[key] = &rldata{
			Count:  1,
			Expire: time.Now().Add(time.Duration(ttlsec) * time.Second),
		}
		return 1, ttlsec, nil
	}
	v.Count++
	rttl := v.Expire.Sub(time.Now()).Seconds()
	if rttl < 1 {
		return v.Count, 0, nil
	}
	return v.Count, int32(rttl), nil
}

// Start starts the internal goroutine that scans the map for expired keys
// and remove them.
func (m *MapLimiter) Start() {
	m.m.Lock()
	defer m.m.Unlock()
	if m.stop != nil {
		return
	}
	m.stop = make(chan struct{})
	ready := make(chan struct{})
	go m.run(ready)
	<-ready
}

// Stop stops the internal goroutine started by Start.
func (m *MapLimiter) Stop() {
	m.m.Lock()
	defer m.m.Unlock()
	if m.stop != nil {
		close(m.stop)
	}
}

func (m *MapLimiter) run(ready chan struct{}) {
	tick := time.NewTicker(m.p)
	close(ready)
	for {
		select {
		case <-m.stop:
			tick.Stop()
			m.m.Lock()
			m.stop = nil
			m.m.Unlock()
		case <-tick.C:
			m.clear()
		}
	}
}

func (m *MapLimiter) clear() {
	now := time.Now()
	m.m.Lock()
	for k, v := range m.s {
		if v.Expire.Sub(now) <= 0 {
			delete(m.s, k)
		}
	}
	m.m.Unlock()
}
