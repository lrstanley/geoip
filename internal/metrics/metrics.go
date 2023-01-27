// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CacheHitCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "geoip_cache_hit_count_total",
			Help: "The total number of cache hits",
		},
		[]string{"bucket"},
	)

	CacheRequestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "geoip_cache_request_count_total",
			Help: "The total number of cache requests (hits + misses)",
		},
		[]string{"bucket"},
	)

	LookupCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "geoip_lookup_total",
			Help: "The total number of lookups for a given lookup type",
		},
		[]string{"type"},
	)
)
